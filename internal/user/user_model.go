package user

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

type RequestPairSchema struct {
	ID       string `json:"id" validate:"required,alphanum"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=16"`
	Pubkey   string `json:"pubkey" validate:"required,base64"`
}

type ResponsePairSchema struct {
	Pubkey string `json:"pubkey"`
}
