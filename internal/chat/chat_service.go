package chat

import (
	"bytes"
	"chat-client/internal/discovery"
	"chat-client/internal/user"
	"chat-client/pkg/encryption"
	"chat-client/pkg/response"
	"chat-client/pkg/store"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/oklog/ulid/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

type ChatService struct {
	ctx              context.Context
	db               *gorm.DB
	s                *store.Store
	discoveryService *discovery.DiscoveryService
}

type IChatService interface {
	CreateChat(input SendMessageSchema) error
	GetMessages(peerId string, cursor uint64) response.Response[[]ChatMessage]
	SendMessage(contact user.ContactModel, input SendMessageSchema) response.Response[string]
	Startup(ctx context.Context)
}

func NewChatService(s *store.Store, db *gorm.DB, discoveryService *discovery.DiscoveryService) *ChatService {
	return &ChatService{s: s, db: db, discoveryService: discoveryService}
}

func (cs *ChatService) CreateChat(input SendMessageSchema) error {
	var contact user.ContactModel

	err := cs.db.First(&contact, "ID = ?", input.Sender).Error
	if err != nil {
		return errors.New("shared key not found")
	}

	if cs.s.Get("key:shared:"+contact.ID) == nil {
		if cs.s.Get("user:password") == nil {
			return errors.New("user password not found")
		}

		contact.SharedKey, err = encryption.PasswordDecrypt([]byte(cs.s.GetString("user:password")), contact.SharedKey)
		if err != nil {
			return errors.New("failed to decrypt shared key")
		}

		cs.s.Set("key:shared:"+contact.ID, contact.SharedKey)
	}

	contact.SharedKey = cs.s.Get("key:shared:" + contact.ID)

	decoded, err := base64.StdEncoding.DecodeString(input.Message)
	if err != nil {
		return errors.New("failed to decode message")
	}

	decrypted, err := encryption.AESDecrypt(contact.SharedKey, decoded)
	if err != nil {
		return errors.New("failed to decrypt message")
	}

	// store message to db
	newMsg := ChatModel{
		ID:      ulid.Now(),
		PeerID:  input.Sender,
		Sender:  input.Sender,
		Message: decoded,
	}
	err = cs.db.Create(&newMsg).Error
	if err != nil {
		return errors.New("db error")
	}

	message := ChatMessage{
		ID:        newMsg.ID,
		Sender:    input.Sender,
		Message:   string(decrypted),
		CreatedAt: newMsg.CreatedAt.Format(time.RFC3339),
	}

	runtime.EventsEmit(cs.ctx, "msg:new:"+message.Sender, message)

	return nil
}

func (cs *ChatService) GetMessages(peerId string, cursor uint64) response.Response[[]ChatMessage] {
	var messages []ChatModel
	var results []ChatMessage

	limit := 3

	if cursor == 0 {
		err := cs.db.Find(&messages, "peer_id = ?", peerId).Limit(limit).Error
		if err != nil {
			return response.New(results).Status(500)
		}
	} else {
		err := cs.db.Find(&messages, "peer_id = ? AND id < ?", peerId, cursor).Limit(limit).Error
		if err != nil {
			return response.New(results).Status(500)
		}
	}

	// check if no older messages
	if len(messages) == 0 {
		return response.New(results).Status(404)
	}

	// retrieve shared key
	var contact user.ContactModel

	err := cs.db.First(&contact, "ID = ?", peerId).Error
	if err != nil {
		return response.New(results).Status(500)
	}

	if cs.s.Get("key:shared:"+contact.ID) == nil {
		if cs.s.Get("user:password") == nil {
			return response.New(results).Status(500)
		}

		contact.SharedKey, err = encryption.PasswordDecrypt([]byte(cs.s.GetString("user:password")), contact.SharedKey)
		if err != nil {
			return response.New(results).Status(500)
		}

		cs.s.Set("key:shared:"+contact.ID, contact.SharedKey)
	}

	contact.SharedKey = cs.s.Get("key:shared:" + contact.ID)

	// decrypt messages
	for _, message := range messages {
		decrypted, err := encryption.AESDecrypt(contact.SharedKey, message.Message)
		if err != nil {
			response.New(results).Status(500)
		}

		results = append(results, ChatMessage{
			ID:        message.ID,
			Sender:    message.Sender,
			Message:   string(decrypted),
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
		})
	}

	return response.New(results)
}

func (cs *ChatService) SendMessage(contact user.ContactModel, input SendMessageSchema) response.Response[ChatMessage] {
	var message ChatMessage
	var err error

	peer := cs.discoveryService.GetPeer(contact.ID)
	if peer.IP == "" {
		return response.New(message).Status(404)
	}

	// retrieve shared key from memory
	contact.SharedKey = cs.s.Get("key:shared:" + contact.ID)
	if contact.SharedKey == nil {
		// get shared key from db if not exist in memory
		err = cs.db.First(&contact, "ID = ?", contact.ID).Error
		if err != nil {
			return response.New(message).Status(500)
		}

		if cs.s.Get("user:password") == nil {
			return response.New(message).Status(500)
		}

		contact.SharedKey, err = encryption.PasswordDecrypt([]byte(cs.s.GetString("user:password")), contact.SharedKey)
		if err != nil {
			return response.New(message).Status(500)
		}

		cs.s.Set("key:shared:"+contact.ID, contact.SharedKey)
	}

	encrypted, err := encryption.AESEncrypt(contact.SharedKey, []byte(input.Message))
	if err != nil {
		return response.New(message).Status(500)
	}

	// encode encrypted message
	encoded := base64.StdEncoding.EncodeToString(encrypted)

	message.Sender = cs.s.GetString("user:id")
	message.Message = encoded

	payload, err := sonic.Marshal(message)
	if err != nil {
		return response.New(message).Status(500)
	}

	url := fmt.Sprintf("http://%s:%d/api/chat/send", peer.IP, discovery.SVC_PORT)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return response.New(message).Status(500)
	}

	if res.StatusCode != http.StatusOK {
		return response.New(message).Status(500)
	}

	// store message to db
	newMsg := ChatModel{
		ID:      ulid.Now(),
		PeerID:  contact.ID,
		Sender:  message.Sender,
		Message: encrypted,
	}
	err = cs.db.Create(&newMsg).Error
	if err != nil {
		return response.New(message).Status(500)
	}

	message.ID = newMsg.ID
	message.Message = input.Message
	message.CreatedAt = newMsg.CreatedAt.Format(time.RFC3339)
	return response.New(message)
}

func (cs *ChatService) Startup(ctx context.Context) {
	cs.ctx = ctx
}
