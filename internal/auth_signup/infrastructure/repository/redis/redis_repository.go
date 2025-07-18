package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/1URose/marketplace/internal/auth_signup/domain/redis/entity"
	"github.com/1URose/marketplace/internal/auth_signup/infrastructure/config/redis"
	"log"
	"time"
)

type Repository struct {
	Client *redis.Client
}

func NewRedisRepository(client *redis.Client) *Repository {
	log.Println("[redis:repo] NewRedisRepository initialized")

	return &Repository{
		Client: client,
	}
}

func (ur *Repository) Set(ctx context.Context, session *entity.Redis) error {
	key := fmt.Sprintf("refresh:%s", session.Email)

	log.Printf("[redis:repo] Set called: key=%q ttl=24h", key)

	data, err := json.Marshal(session.RefreshToken)

	if err != nil {

		log.Printf("[redis:repo][ERROR] marshal token failed: %v", err)

		return err
	}

	if err = ur.Client.Connection.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {

		log.Printf("[redis:repo][ERROR] SET command failed for key=%q: %v", key, err)

		return err
	}

	log.Printf("[redis:repo] Set succeeded for key=%q", key)

	return nil
}

func (ur *Repository) Get(ctx context.Context, email string) (*entity.Redis, error) {
	key := fmt.Sprintf("refresh:%s", email)

	log.Printf("[redis:repo] Get called: key=%q", key)

	data, err := ur.Client.Connection.Get(ctx, key).Result()

	if err != nil {

		log.Printf("[redis:repo][ERROR] GET command failed for key=%q: %v", key, err)

		return nil, err
	}

	var session entity.Redis

	if err = json.Unmarshal([]byte(data), &session); err != nil {

		log.Printf("[redis:repo][ERROR] unmarshal data failed for key=%q: %v", key, err)

		return nil, err
	}

	log.Printf("[redis:repo] Get succeeded for key=%q: token=%q", key, session.RefreshToken)

	return &session, nil
}
