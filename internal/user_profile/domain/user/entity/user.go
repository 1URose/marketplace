package entity

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:"password_hash,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

func NewUser(email, passwordHash string) *User {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
	}
}
