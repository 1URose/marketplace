package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	pool *pgxpool.Pool
}

func NewClient() (*Client, error) {

	log.Println("[postgresql:user] loading config from env")

	cfg := ReadConfigFromEnv()

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB,
	)

	log.Printf("[postgresql:user] connecting to %s", connStr)

	pgConfig, err := pgxpool.ParseConfig(connStr)

	if err != nil {

		log.Printf("[postgresql:user][ERROR] parse config: %v", err)

		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	pgConfig.MaxConns = cfg.MaxConns
	pgConfig.MinConns = cfg.MinConns
	pgConfig.MaxConnLifetime = cfg.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)

	if err != nil {

		log.Printf("[postgresql:user][ERROR] create pool: %v", err)

		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	log.Println("[postgresql:user] connection pool established")

	if err = pool.Ping(context.Background()); err != nil {
		log.Printf("[postgresql:user][ERROR] ping database failed: %v", err)

		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("[postgresql:user] database ping successful")

	return &Client{pool: pool}, nil
}

func (c *Client) Close() {

	log.Println("[postgresql:user] closing connection pool")

	c.pool.Close()

	log.Println("[postgresql:user] connection pool closed")
}

func (c *Client) GetPool() *pgxpool.Pool {
	log.Println("[postgresql:user] GetPool called, returning connection pool")
	return c.pool
}
