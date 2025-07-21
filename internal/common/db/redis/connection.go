package redis

import (
	"context"
	"fmt"
	redisCfd "github.com/1URose/marketplace/internal/common/config/redis"
	"log"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	Connection *redis.Client
}

func NewRedisClient(cfg *redisCfd.Config) (*Client, error) {

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	log.Printf("[redis] connecting to %s (user=%s db=%d)", addr, cfg.User, cfg.DB)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {

		log.Printf("[redis][ERROR] ping failed: %v", err)

		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("[redis][ERROR] ping failed: %v", err)
		if cerr := client.Close(); cerr != nil {
			log.Printf("[redis][ERROR] failed to close client after ping error: %v", cerr)
		}
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	log.Println("[redis] connection established")

	return &Client{Connection: client}, nil
}

func (r *Client) Close() error {
	log.Println("[redis] closing connection")

	if err := r.Connection.Close(); err != nil {

		log.Printf("[redis][ERROR] failed to close Redis connection: %v", err)

		return fmt.Errorf("failed to close Redis connection: %w", err)
	}

	log.Println("[redis] connection closed")

	return nil
}
