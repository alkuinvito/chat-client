package user

import (
	"bytes"
	"chat-client/internal/discovery"
	"chat-client/pkg/encryption"
	"chat-client/pkg/response"
	"chat-client/pkg/store"
	"context"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	ctx              context.Context
	db               *gorm.DB
	discoveryService *discovery.DiscoveryService
	s                *store.Store
	router           *fiber.App
}

type IUserService interface {
	GeneratePairingCode() response.Response[string]
	generateSharedKey() error
	GetContacts() response.Response[[]ContactModel]
	getDefaultUser() (UserModel, error)
	GetProfile() response.Response[UserProfile]
	HandleUserPairing(input RequestPairSchema) (ResponsePairSchema, error)
	loadPrivateKey(password []byte) (*ecdh.PrivateKey, error)
	Login(username, password string) response.Response[UserProfile]
	Register(username, password string) response.Response[UserProfile]
	RequestPairing(input RequestPairSchema) response.Response[string]
	ScanPeers() response.Response[[]discovery.PeerModel]
	Startup(ctx context.Context)
}

func NewUserService(s *store.Store, db *gorm.DB, router *fiber.App, discoveryService *discovery.DiscoveryService) *UserService {
	return &UserService{
		s:                s,
		db:               db,
		discoveryService: discoveryService,
		router:           router,
	}
}

// Generate 6-digit pairing code and expires in 60 seconds
func (us *UserService) GeneratePairingCode() response.Response[string] {
	nA, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return response.New("failed to generate pairing code").Status(500)
	}

	nB, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return response.New("failed to generate pairing code").Status(500)
	}

	nC, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		return response.New("failed to generate pairing code").Status(500)
	}

	pairingCode := fmt.Sprintf("%02d%02d%02d", nA.Int64(), nB.Int64(), nC.Int64())

	// pairing code expires in 60 seconds
	us.s.SetEx("pair:code", []byte(pairingCode), time.Second*60)

	return response.New(pairingCode)
}

// Generate shared key from remote public key
func (us *UserService) generateSharedKey(remotePubkey []byte) ([]byte, []byte, error) {
	if us.s.Get("user:password") == nil {
		return nil, nil, errors.New("user password not found")
	}

	remote, err := ecdh.P256().NewPublicKey(remotePubkey)
	if err != nil {
		return nil, nil, errors.New("invalid remote public key")
	}

	priv, err := us.loadPrivateKey([]byte(us.s.GetString("user:password")))
	if err != nil {
		return nil, nil, errors.New("failed to load private key")
	}

	shared, err := encryption.GenerateSharedKey(priv, remote)
	if err != nil {
		return nil, nil, errors.New("failed to generate shared key")
	}

	sharedEnc, err := encryption.PasswordEncrypt([]byte(us.s.GetString("user:password")), shared)
	if err != nil {
		return nil, nil, errors.New("failed to encrypt shared key")
	}

	return shared, sharedEnc, nil
}

// Get all contacts
func (us *UserService) GetContacts() response.Response[[]ContactModel] {
	var result []ContactModel

	err := us.db.Find(&result).Error
	if err != nil {
		return response.New(result).Status(500)
	}

	return response.New(result)
}

// Return created default user
func (us *UserService) getDefaultUser() (UserModel, error) {
	var result UserModel

	err := us.db.First(&result).Error
	if err != nil {
		return result, errors.New("user not found")
	}

	return result, nil
}

// Get user profile
func (us *UserService) GetProfile() response.Response[UserProfile] {
	var result UserModel

	err := us.db.First(&result).Error
	if err != nil {
		return response.New(result.toProfile()).Status(404)
	}

	return response.New(result.toProfile())
}

// Handle user pairing and create shared key
func (us *UserService) HandleUserPairing(input InitPairSchema) (ResponsePairSchema, error) {
	var result ResponsePairSchema

	pairCode := us.s.GetString("pair:code")
	if pairCode == "" {
		return result, errors.New("pairing code not found")
	}

	// generate sha256sum of pair code
	hash := sha256.Sum256([]byte(pairCode))
	hashString := hex.EncodeToString([]byte(hash[:]))

	// checksum of both paircode
	if input.Code != hashString {
		return result, errors.New("pairing code incorrect")
	}

	// check for existing contact
	var oldContact ContactModel
	err := us.db.First(&oldContact, "ID = ?", input.ID).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return result, errors.New("db error")
		}
	}

	if oldContact.ID == input.ID {
		return result, errors.New("user already paired")
	}

	if us.s.Get("key:public") == nil {
		return result, errors.New("public key not found")
	}

	// decode base64 pubkey to bytes
	decoded, err := base64.StdEncoding.DecodeString(input.Pubkey)
	if err != nil {
		return result, errors.New("invalid base64 pubkey")
	}

	// decrypt pubkey using pre-shared passcode
	decrypted, err := encryption.PasswordDecrypt([]byte(pairCode), decoded)
	if err != nil {
		return result, errors.New("invalid encrypted pubkey")
	}

	shared, sharedEnc, err := us.generateSharedKey(decrypted)
	if err != nil {
		if err.Error() == "invalid remote public key" {
			return result, err
		}

		return result, errors.New("failed to generate shared key")
	}

	contact := ContactModel{
		ID:        input.ID,
		Username:  input.Username,
		SharedKey: sharedEnc,
	}

	// save the newly paired contact
	err = us.db.Create(&contact).Error
	if err != nil {
		return result, err
	}

	// store the shared key in memory
	us.s.Set("key:shared:"+input.ID, shared)

	// broadcast for new contact
	contact.SharedKey = nil
	runtime.EventsEmit(us.ctx, "pair:new", contact)

	// encrypt public key using pre-shared passcode
	encrypted, err := encryption.PasswordEncrypt([]byte(pairCode), us.s.Get("key:public"))

	result.Pubkey = base64.StdEncoding.EncodeToString(encrypted)

	return result, nil
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

	// check for username
	if result.Username != username {
		return response.New(result.toProfile()).Status(401)
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

	// start service query
	go us.discoveryService.QueryService()

	// start chat server
	go us.router.Listen(fmt.Sprintf(":%d", discovery.SVC_PORT))

	return response.New(result.toProfile())
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

func (us *UserService) RequestPairing(input RequestPairSchema) response.Response[string] {
	userId := us.s.GetString("user:id")
	if userId == "" {
		return response.New("userId not found").Status(500)
	}

	username := us.s.GetString("user:username")
	if username == "" {
		return response.New("username not found").Status(500)
	}

	if us.s.Get("key:public") == nil {
		return response.New("public key not found").Status(500)
	}

	// encrypt ecdh pubkey using pre-shared passcode
	encrypted, err := encryption.PasswordEncrypt([]byte(input.Code), us.s.Get("key:public"))
	if err != nil {
		return response.New("failed to encrypt public key")
	}

	// encode pubkey to base64
	encoded := base64.StdEncoding.EncodeToString(encrypted)

	// hash passcode using sha256
	hash := sha256.Sum256([]byte(input.Code))
	hashString := hex.EncodeToString(hash[:])

	initReq := InitPairSchema{
		ID:       userId,
		Username: username,
		Pubkey:   encoded,
		Code:     hashString,
	}

	payload, err := sonic.Marshal(&initReq)
	if err != nil {
		return response.New("failed to generate json").Status(500)
	}

	peer := us.discoveryService.GetPeer(input.ID)
	if peer.IP == "" {
		return response.New("peer is not found")
	}

	url := fmt.Sprintf("http://%s:%d/api/user/pair", peer.IP, discovery.SVC_PORT)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return response.New("failed to initiate pair").Status(500)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response.New("failed to read response body").Status(500)
	}

	if res.StatusCode != http.StatusOK {
		var resErr response.ErrorResponseSchema
		err = sonic.Unmarshal(body, &resErr)
		if err != nil {
			return response.New("failed to read response error schema").Status(500)
		}

		return response.New(resErr.Error).Status(res.StatusCode)
	}

	var resPair ResponsePairSchema
	err = sonic.Unmarshal(body, &resPair)
	if err != nil {
		return response.New("failed to read response pair schema").Status(500)
	}

	// decode base64 pubkey to bytes
	decoded, err := base64.StdEncoding.DecodeString(resPair.Pubkey)
	if err != nil {
		return response.New("invalid base64 pubkey").Status(500)
	}

	// decrypt pubkey using pre-shared passcode
	decrypted, err := encryption.PasswordDecrypt([]byte(input.Code), decoded)
	if err != nil {
		return response.New("invalid encrypted pubkey").Status(500)
	}

	shared, sharedEnc, err := us.generateSharedKey(decrypted)
	if err != nil {
		if err.Error() == "invalid remote public key" {
			return response.New(err.Error()).Status(500)
		}

		return response.New("failed to generate shared key").Status(500)
	}

	contact := ContactModel{
		ID:        input.ID,
		Username:  input.Username,
		SharedKey: sharedEnc,
	}

	// save the newly paired contact
	err = us.db.Create(&contact).Error
	if err != nil {
		return response.New("failed to store contact").Status(500)
	}

	// store the shared key in memory
	us.s.Set("key:shared:"+input.ID, shared)

	// broadcast for new contact
	contact.SharedKey = nil
	runtime.EventsEmit(us.ctx, "pair:new", contact)

	return response.New("paired successfully")
}

func (us *UserService) ScanPeers() response.Response[[]discovery.PeerModel] {
	var result []discovery.PeerModel
	isContact := make(map[string]bool)

	contacts := us.GetContacts()
	if contacts.Code != 200 {
		return response.New(result).Status(500)
	}

	for _, contact := range contacts.Data {
		isContact[contact.ID] = true
	}

	// force refresh query
	us.discoveryService.RefreshQuery()

	// get query results
	peers := us.discoveryService.GetPeers()
	if peers.Code != 200 {
		return response.New(result).Status(500)
	}

	// ignore active peers that is in contact list
	for _, peer := range peers.Data {
		if isContact[peer.ID] {
			continue
		}

		result = append(result, peer)
	}

	return response.New(result)
}

func (us *UserService) Startup(ctx context.Context) {
	us.ctx = ctx
}
