package db

import (
	"fmt"
	"github.com/1URose/marketplace/internal/common/config"
	"github.com/1URose/marketplace/internal/common/db/postgresql"
	"github.com/1URose/marketplace/internal/common/db/redis"
	"log"
)

type Connections struct {
	PostgresConn *postgresql.Client
	RedisConn    *redis.Client
}

func NewConnections(cfg *config.GeneralConfig) (*Connections, error) {
	log.Println("Creating user Postgres connection")

	userPostgresConn, err := postgresql.NewClient(cfg.PostgresConfig)

	if err != nil {
		log.Printf("ERROR: user Postgres connection failed: %v", err)

		return nil, fmt.Errorf("failed to create user postgres connection: %w", err)
	}

	log.Println("User Postgres connected")

	log.Println("Creating Redis connection")

	redisConn, err := redis.NewRedisClient(cfg.RedisConfig)
	if err != nil {
		log.Printf("ERROR: Redis connection failed: %v", err)

		return nil, fmt.Errorf("failed to create Redis connection: %w", err)
	}

	log.Println("Redis connected")

	return &Connections{
		PostgresConn: userPostgresConn,
		RedisConn:    redisConn,
	}, nil
}

func (c *Connections) Close() {
	log.Println("Closing database connections…")

	if c.PostgresConn != nil {
		c.PostgresConn.Close()

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
