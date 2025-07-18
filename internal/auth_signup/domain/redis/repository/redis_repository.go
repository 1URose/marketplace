package repository

import (
	"context"
	"github.com/1URose/marketplace/internal/auth_signup/domain/redis/entity"
)

type RedisRepository interface {
	Set(ctx context.Context, hash *entity.Redis) error
	Get(ctx context.Context, email string) (*entity.Redis, error)
}
