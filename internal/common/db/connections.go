package db

import (
	"fmt"
	"github.com/1URose/marketplace/internal/auth_signup/infrastructure/config/redis"
	"github.com/1URose/marketplace/internal/common/db/postgresql"
	"log"
)

type Connections struct {
	UserPostgresConn *postgresql.Client
	RedisConn        *redis.Client
}

func NewConnections() (*Connections, error) {
	log.Println("Creating user Postgres connection")

	userPostgresConn, err := postgresql.NewClient()

	if err != nil {
		log.Printf("ERROR: user Postgres connection failed: %v", err)

		return nil, fmt.Errorf("failed to create user postgres connection: %w", err)
	}

	log.Println("User Postgres connected")

	log.Println("Creating Redis connection")

	redisConn, err := redis.NewRedisClient()
	if err != nil {
		log.Printf("ERROR: Redis connection failed: %v", err)

		return nil, fmt.Errorf("failed to create Redis connection: %w", err)
	}

	log.Println("Redis connected")

	return &Connections{
		UserPostgresConn: userPostgresConn,
		RedisConn:        redisConn,
	}, nil
}

func (c *Connections) Close() {
	log.Println("Closing database connections…")

	if c.UserPostgresConn != nil {
		c.UserPostgresConn.Close()

		log.Println("User Postgres closed")
	}

	if c.RedisConn != nil {
		if err := c.RedisConn.Close(); err != nil {

			log.Printf("ERROR: failed to close Redis: %v", err)

		} else {

			log.Println("Redis closed")

		}
	}
}
