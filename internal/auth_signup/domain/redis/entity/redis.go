package entity

type Redis struct {
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
}

func NewSession(email string, refreshToken string) *Redis {
	return &Redis{
		Email:        email,
		RefreshToken: refreshToken,
	}
}
