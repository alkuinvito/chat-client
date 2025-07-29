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

	"github.com/bytedance/sonic"
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
	CreateChat(payload ChatMessage) error
	SendMessage(contact user.ContactModel, message ChatMessage) response.Response[string]
	Startup(ctx context.Context)
}

func NewChatService(s *store.Store, db *gorm.DB, discoveryService *discovery.DiscoveryService) *ChatService {
	return &ChatService{s: s, db: db, discoveryService: discoveryService}
}

func (cs *ChatService) CreateChat(payload ChatMessage) error {
	var contact user.ContactModel

	err := cs.db.First(&contact, "ID = ?", payload.Sender).Error
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

	decoded, err := base64.StdEncoding.DecodeString(payload.Message)
	if err != nil {
		return errors.New("failed to decode message")
	}

	decrypted, err := encryption.AESDecrypt(contact.SharedKey, decoded)
	if err != nil {
		return errors.New("failed to decrypt message")
	}

	payload.Message = string(decrypted)

	runtime.EventsEmit(cs.ctx, "msg:new:"+payload.Sender, payload)

	return nil
}

func (cs *ChatService) SendMessage(contact user.ContactModel, message ChatMessage) response.Response[string] {
	var err error

	peer := cs.discoveryService.GetPeer(contact.ID)
	if peer.IP == "" {
		return response.New("peer not found").Status(404)
	}

	// retrieve shared key from memory
	contact.SharedKey = cs.s.Get("key:shared:" + contact.ID)
	if contact.SharedKey == nil {
		// get shared key from db if not exist in memory
		err = cs.db.First(&contact, "ID = ?", contact.ID).Error
		if err != nil {
			return response.New("shared key not found").Status(500)
		}

		if cs.s.Get("user:password") == nil {
			return response.New("user password not found").Status(500)
		}

		contact.SharedKey, err = encryption.PasswordDecrypt([]byte(cs.s.GetString("user:password")), contact.SharedKey)
		if err != nil {
			return response.New("failed to decrypt shared key").Status(500)
		}

		cs.s.Set("key:shared:"+contact.ID, contact.SharedKey)
	}

	encrypted, err := encryption.AESEncrypt(contact.SharedKey, []byte(message.Message))
	if err != nil {
		return response.New("failed to encrypt message").Status(500)
	}

	encoded := base64.StdEncoding.EncodeToString(encrypted)
	message.Sender = cs.s.GetString("user:id")
	message.Message = encoded

	payload, err := sonic.Marshal(message)
	if err != nil {
		return response.New("failed to send message").Status(500)
	}

	url := fmt.Sprintf("http://%s:%d/api/chat/send", peer.IP, discovery.SVC_PORT)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return response.New("failed to send message").Status(500)
	}

	if res.StatusCode != http.StatusOK {
		return response.New("failed to send message").Status(500)
	}

	return response.New("message sent")
}

func (cs *ChatService) Startup(ctx context.Context) {
	cs.ctx = ctx
}
