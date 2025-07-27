package user

import (
	"chat-client/internal/chat"
	"chat-client/internal/discovery"
	"chat-client/internal/router"
	"chat-client/pkg/encryption"
	"chat-client/pkg/response"
	"chat-client/pkg/store"
	"context"
	"log"

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
	GetDefaultUser() response.Response[UserProfile]
	GetProfile() response.Response[UserProfile]
	Login(username, password string) response.Response[UserProfile]
	Register(username, password string) response.Response[UserProfile]
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

func (as *UserService) GetDefaultUser() response.Response[UserProfile] {
	var result UserModel

	err := as.db.First(&result).Error
	if err != nil {
		return response.New(result.toProfile()).Status(404)
	}

	return response.New(result.toProfile())
}

func (as *UserService) GetProfile() response.Response[UserProfile] {
	var result UserModel

	username, err := as.s.GetString("username")
	if err != nil {
		return response.New(result.toProfile()).Status(404)
	}

	result.Username = username
	return response.New(result.toProfile())
}

func (as *UserService) Login(username, password string) response.Response[UserProfile] {
	var result UserModel

	err := as.db.First(&result).Error
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
	as.s.Set("username", []byte(username))

	// store pubkey in memory
	as.s.Set("key:public", pubkey)

	// start broadcasting the service
	go as.discoveryService.BroadcastService(username)

	// start chat server
	go as.router.Handle()

	return response.New(result.toProfile())
}

func (as *UserService) Register(username, password string) response.Response[UserProfile] {
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

	err = as.db.Create(&user).Error
	if err != nil {
		log.Println(err)
		return response.New(user.toProfile()).Status(500)
	}

	// store username in memory
	as.s.Set("username", []byte(username))

	// store public key in memory
	as.s.Set("key:public", priv.PublicKey().Bytes())

	// start broadcasting the service
	go as.discoveryService.BroadcastService(username)

	// start chat server
	go as.router.Handle()

	return response.New(user.toProfile())
}

func (as *UserService) Startup(ctx context.Context) {
	as.ctx = ctx
}
