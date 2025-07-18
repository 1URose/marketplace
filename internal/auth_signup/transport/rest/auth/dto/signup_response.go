package dto

import (
	"github.com/1URose/marketplace/internal/user_profile/domain/user/entity"
	"github.com/1URose/marketplace/internal/user_profile/transport/rest/user/dto"
)

type SignUpResponse struct {
	User *dto.UserResponse
}

func NewSignUpResponse(user *entity.User) *SignUpResponse {
	return &SignUpResponse{
		User: dto.NewUserResponse(user),
	}
}
