package auth

type UserModel struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"not null" validate:"alphanum,min=3,max=16"`
	Password string `json:"password" gorm:"not null" validate:"min=8,max=32"`
}

type UserProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (um *UserModel) toProfile() UserProfile {
	return UserProfile{ID: um.ID, Username: um.Username}
}
