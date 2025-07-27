package user

import (
	"bytes"
	"chat-client/internal/chat"
	"chat-client/internal/discovery"
	"chat-client/internal/router"
	"chat-client/pkg/encryption"
	"chat-client/pkg/response"
	"chat-client/pkg/store"
	"context"
	"crypto/ecdh"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	chatService      *chat.ChatService
	ctx              context.Context
	db               *gorm.DB
	discoveryService *discovery.DiscoveryService
	s                *store.Store
	router           *router.Router
}

type IUserService interface {
	generateSharedKey() error
	getDefaultUser() (UserModel, error)
	GetProfile() response.Response[UserProfile]
	loadPrivateKey(password []byte) (*ecdh.PrivateKey, error)
	Login(username, password string) response.Response[UserProfile]
	PairUser(input RequestPairSchema) (ResponsePairSchema, error)
	Register(username, password string) response.Response[UserProfile]
	RequestPair(peer discovery.PeerModel) response.Response[string]
	Startup(ctx context.Context)
}

func NewUserService(s *store.Store, db *gorm.DB, router *router.Router, discoveryService *discovery.DiscoveryService, chatService *chat.ChatService) *UserService {
	return &UserService{
		s:                s,
		db:               db,
		discoveryService: discoveryService,
		chatService:      chatService,
		router:           router,
	}
}

func (us *UserService) generateSharedKey(remotePubkey string) ([]byte, []byte, error) {
	password, err := us.s.Get("user:password")
	if err != nil {
		return nil, nil, errors.New("user password not found")
	}

	encoded, err := base64.StdEncoding.DecodeString(remotePubkey)
	if err != nil {
		return nil, nil, errors.New("invalid remote public key")
	}

	remote, err := ecdh.P256().NewPublicKey(encoded)
	if err != nil {
		return nil, nil, errors.New("invalid remote public key")
	}

	priv, err := us.loadPrivateKey(password)
	if err != nil {
		return nil, nil, errors.New("failed to load private key")
	}

	shared, err := encryption.GenerateSharedKey(priv, remote)
	if err != nil {
		return nil, nil, errors.New("failed to generate shared key")
	}

	sharedEnc, err := encryption.PasswordEncrypt(password, shared)
	if err != nil {
		return nil, nil, errors.New("failed to encrypt shared key")
	}

	return shared, sharedEnc, nil
}

func (us *UserService) getDefaultUser() (UserModel, error) {
	var result UserModel

	err := us.db.First(&result).Error
	if err != nil {
		return result, errors.New("user not found")
	}

	return result, nil
}

func (us *UserService) GetProfile() response.Response[UserProfile] {
	var result UserModel

	err := us.db.First(&result).Error
	if err != nil {
		return response.New(result.toProfile()).Status(404)
	}

	return response.New(result.toProfile())
}

func (us *UserService) loadPrivateKey(password []byte) (*ecdh.PrivateKey, error) {
	user, err := us.getDefaultUser()
	if err != nil {
		return nil, err
	}

	decrypted, err := encryption.PasswordDecrypt(password, user.PrivKey)
	if err != nil {
		return nil, err
	}

	priv, err := ecdh.P256().NewPrivateKey(decrypted)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func (us *UserService) Login(username, password string) response.Response[UserProfile] {
	var result UserModel

	err := us.db.First(&result).Error
	if err != nil {
		return response.New(result.toProfile()).Status(404)
	}

	// check for password
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))
	if err != nil {
		return response.New(result.toProfile()).Status(401)
	}

	// get public key
	pubkey, err := encryption.PasswordDecrypt([]byte(password), result.PubKey)
	if err != nil {
		return response.New(result.toProfile()).Status(500)
	}

	// store username in memory
	us.s.Set("user:username", []byte(username))

	// store password in memory
	us.s.Set("user:id", []byte(result.ID))

	// store password in memory
	us.s.Set("user:password", []byte(password))

	// store pubkey in memory
	us.s.Set("key:public", pubkey)

	// start broadcasting the service
	go us.discoveryService.BroadcastService(result.ID, username)

	// start chat server
	go us.router.Handle()

	return response.New(result.toProfile())
}

func (us *UserService) PairUser(input RequestPairSchema) (ResponsePairSchema, error) {
	var result ResponsePairSchema
	pubkey, err := us.s.Get("key:public")
	if err != nil {
		return result, errors.New("public key not found")
	}

	shared, sharedEnc, err := us.generateSharedKey(input.Pubkey)
	if err != nil {
		if err.Error() == "invalid remote public key" {
			return result, err
		}

		return result, errors.New("unknown error")
	}

	contact := ContactModel{
		ID:        input.ID,
		Username:  input.Username,
		SharedKey: sharedEnc,
	}

	err = us.db.Create(&contact).Error
	if err != nil {
		return result, errors.New("failed to store contact to db")
	}

	us.s.Set("key:shared:"+input.ID, shared)

	result.Pubkey = base64.StdEncoding.EncodeToString(pubkey)

	return result, nil

}

func (us *UserService) Register(username, password string) response.Response[UserProfile] {
	var user UserModel

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return response.New(user.toProfile()).Status(500)
	}

	priv, err := encryption.GeneratePrivateKey()
	if err != nil {
		log.Println(err)
		return response.New(user.toProfile()).Status(500)
	}

	privEnc, err := encryption.PasswordEncrypt([]byte(password), priv.Bytes())
	if err != nil {
		log.Println(err)
		return response.New(user.toProfile()).Status(500)
	}

	pubEnc, err := encryption.PasswordEncrypt([]byte(password), priv.PublicKey().Bytes())
	if err != nil {
		log.Println(err)
		return response.New(user.toProfile()).Status(500)
	}

	user = UserModel{
		ID:       ulid.Make().String(),
		Username: username,
		Password: string(hashed),
		PrivKey:  privEnc,
		PubKey:   pubEnc,
	}

	err = us.db.Create(&user).Error
	if err != nil {
		log.Println(err)
		return response.New(user.toProfile()).Status(500)
	}

	return response.New(user.toProfile())
}

func (us *UserService) RequestPair(peer discovery.PeerModel) response.Response[string] {
	pubkey, err := us.s.Get("key:public")
	if err != nil {
		return response.New("public key not found").Status(500)
	}

	userId, err := us.s.Get("user:id")
	if err != nil {
		return response.New("userId not found").Status(500)
	}

	username, err := us.s.Get("user:username")
	if err != nil {
		return response.New("username not found").Status(500)
	}

	pubkeyEnc := base64.StdEncoding.EncodeToString(pubkey)

	req := RequestPairSchema{
		ID:       string(userId),
		Username: string(username),
		Pubkey:   pubkeyEnc,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return response.New("failed to produce json").Status(500)
	}

	url := fmt.Sprintf("http://%s:%d/api/user/pair", peer.IP, discovery.SVC_PORT)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return response.New("failed to send pair request").Status(500)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response.New("failed to read response").Status(500)
	}

	log.Println(res.StatusCode, string(body))

	remotePubkey := ResponsePairSchema{}
	err = json.Unmarshal(body, &remotePubkey)
	if err != nil {
		return response.New("failed to read response").Status(500)
	}

	shared, sharedEnc, err := us.generateSharedKey(remotePubkey.Pubkey)
	if err != nil {
		log.Println(err)
		return response.New("failed to generate shared key").Status(500)
	}

	contact := ContactModel{
		ID:        peer.ID,
		Username:  peer.Username,
		SharedKey: sharedEnc,
	}

	err = us.db.Create(&contact).Error
	if err != nil {
		log.Println(err)
		return response.New("failed to save new contact").Status(500)
	}

	us.s.Set("key:shared:"+peer.ID, shared)

	return response.New("successfully paired")
}

func (us *UserService) Startup(ctx context.Context) {
	us.ctx = ctx
}
