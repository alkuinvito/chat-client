package auth

import (
	"chat-client/internal/chat"
	"chat-client/internal/discovery"
	"chat-client/pkg/response"
	"chat-client/pkg/store"
	"context"
	"log"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	chatService      *chat.ChatService
	ctx              context.Context
	db               *gorm.DB
	discoveryService *discovery.DiscoveryService
	s                *store.Store
}

type IAuthService interface {
	GetDefaultUser() response.Response[UserProfile]
	GetProfile() response.Response[UserProfile]
	Login(username, password string) response.Response[UserProfile]
	Register(username, password string) response.Response[UserProfile]
	Startup(ctx context.Context)
}

func NewAuthService(s *store.Store, db *gorm.DB, discoveryService *discovery.DiscoveryService, chatService *chat.ChatService) *AuthService {
	return &AuthService{
		s:                s,
		db:               db,
		discoveryService: discoveryService,
		chatService:      chatService,
	}
}

func (as *AuthService) GetDefaultUser() response.Response[UserProfile] {
	var result UserModel

	err := as.db.First(&result).Error
	if err != nil {
		return response.New(result.toProfile()).Status(404)
	}

	result.Password = ""
	return response.New(result.toProfile())
}

func (as *AuthService) GetProfile() response.Response[UserProfile] {
	var result UserModel

	username, err := as.s.Get("username")
	if err != nil {
		return response.New(result.toProfile()).Status(404)
	}

	result.Username = username
	return response.New(result.toProfile())

}

func (as *AuthService) Login(username, password string) response.Response[UserProfile] {
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

	// store username in-memory
	as.s.Set("username", username)

	return response.New(result.toProfile())
}

func (as *AuthService) Register(username, password string) response.Response[UserProfile] {
	var user UserModel

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return response.New(user.toProfile()).Status(500)
	}

	user = UserModel{
		ID:       ulid.Make().String(),
		Username: username,
		Password: string(hashed),
	}

	err = as.db.Create(&user).Error
	if err != nil {
		log.Println(err)
		return response.New(user.toProfile()).Status(500)
	}

	// store username in-memory
	as.s.Set("username", username)

	// start broadcasting the service
	go as.discoveryService.BroadcastService(username)

	// start chat server
	go as.chatService.ServeChat()

	return response.New(user.toProfile())
}

func (as *AuthService) Startup(ctx context.Context) {
	as.ctx = ctx
}
