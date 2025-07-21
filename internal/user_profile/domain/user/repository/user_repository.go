package repository

import (
	"context"
	"github.com/1URose/marketplace/internal/user_profile/domain/user/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}
