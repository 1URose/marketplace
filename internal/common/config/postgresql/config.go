package postgresql

import (
	"github.com/1URose/marketplace/internal/common/settings"
	"log"
	"time"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	DB              string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
}

func NewConfig(
	host, port, user, password, DB string,
	maxConns, minConns int32,
	maxConnLifetime time.Duration,
) *Config {
	return &Config{
		Host:            host,
		Port:            port,
		User:            user,
		Password:        password,
		DB:              DB,
		MaxConns:        maxConns,
		MinConns:        minConns,
		MaxConnLifetime: maxConnLifetime,
	}
}

func LoadPGConfigFromEnv() *Config {
	log.Println("[postgresql:config] reading config from env")

	const (
		envHost        = "PG_HOST"
		envPort        = "PG_PORT"
		envUser        = "PG_USER"
		envPassword    = "PG_PASSWORD"
		envDB          = "PG_DB"
		envMaxConns    = "PG_MAX_CONNS"
		envMinConns    = "PG_MIN_CONNS"
		envMaxLifetime = "PG_MAX_CONN_LIFETIME"
	)

	host := settings.GetEnvSrt(envHost)
	port := settings.GetEnvSrt(envPort)
	user := settings.GetEnvSrt(envUser)
	pass := settings.GetEnvSrt(envPassword)
	db := settings.GetEnvSrt(envDB)

	maxConns, err := settings.GetEnvInt32(envMaxConns)
	if err != nil {
		log.Panicf("[postgresql:config][FATAL] invalid %s: %v", envMaxConns, err)
	}

	minConns, err := settings.GetEnvInt32(envMinConns)
	if err != nil {
		log.Panicf("[postgresql:config][FATAL] invalid %s: %v", envMinConns, err)
	}

	maxLifetimeSec, err := settings.GetEnvInt32(envMaxLifetime)
	if err != nil {
		log.Panicf("[postgresql:config][FATAL] invalid %s: %v", envMaxLifetime, err)
	}

	cfg := NewConfig(
		host, port, user, pass, db,
		maxConns, minConns,
		time.Duration(maxLifetimeSec)*time.Second,
	)

	log.Printf(
		"[postgresql:config] loaded: host=%s port=%s user=%s db=%s max_conns=%d min_conns=%d max_conn_lifetime=%s",
		host, port, user, db, maxConns, minConns, maxLifetimeSec,
	)

	return cfg
}
