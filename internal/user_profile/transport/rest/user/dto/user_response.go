package dto

import (
	"github.com/1URose/marketplace/internal/user_profile/domain/user/entity"
	"time"
)

type UserResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func NewUserResponse(user *entity.User) *UserResponse {
	createdAt := user.CreatedAt.Format(time.RFC3339)
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: createdAt,
	}
}
