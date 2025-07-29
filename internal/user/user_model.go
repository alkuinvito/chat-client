package user

import "chat-client/internal/discovery"

type UserModel struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"not null" validate:"required,alphanum,min=3,max=16"`
	Password string `json:"password" gorm:"not null" validate:"required,min=8,max=32"`
	PrivKey  []byte `gorm:"not null"`
	PubKey   []byte `gorm:"not null"`
}

type UserProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (um *UserModel) toProfile() UserProfile {
	return UserProfile{ID: um.ID, Username: um.Username}
}

type ContactModel struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"not null"`
	SharedKey []byte `gorm:"not null"`
}

type PairRequestModel struct {
	ID       string `json:"id" gorm:"primaryKey" validate:"required,alphanum"`
	Username string `json:"username" gorm:"not null" validate:"required,alphanum,min=3,max=16"`
	Code     string `gorm:"not null"`
	Type     string `json:"type" gorm:"not null"`
}

type InitPairingRequest struct {
	Peer discovery.PeerModel
	Code string
}

type InitPairSchema struct {
	ID       string `json:"id" validate:"required,alphanum"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=16"`
	Code     string `json:"code" validate:"required,numeric,length=4"`
	Pubkey   string `json:"pubkey" validate:"required,base64"`
}

type RequestPairSchema struct {
	ID       string `json:"id" validate:"required,alphanum"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=16"`
	Code     string `json:"code" validate:"required,numeric,length=6"`
}

type ResponsePairSchema struct {
	Pubkey string `json:"pubkey"`
}
