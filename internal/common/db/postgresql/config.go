package postgresql

import (
	"log"
	"os"
	"strconv"
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

func ReadConfigFromEnv() *Config {
	log.Println("[postgresql:user] reading config from env")

	host := getEnv("USER_PG_HOST")
	port := getEnv("USER_PG_PORT")
	user := getEnv("USER_PG_USER")
	pass := getEnv("USER_PG_PASSWORD")
	db := getEnv("USER_PG_DB")

	maxConns, err := parseEnvInt32("USER_PG_MAX_CONNS")

	if err != nil {
		log.Panicf("[postgresql:user][FATAL] invalid USER_PG_MAX_CONNS: %v", err)
	}

	minConns, err := parseEnvInt32("USER_PG_MIN_CONNS")

	if err != nil {
		log.Panicf("[postgresql:user][FATAL] invalid USER_PG_MIN_CONNS: %v", err)
	}

	maxLifetimeSec, err := parseEnvInt32("USER_PG_MAX_CONN_LIFETIME")

	if err != nil {
		log.Panicf("[postgresql:user][FATAL] invalid USER_PG_MAX_CONN_LIFETIME: %v", err)
	}

	cfg := NewConfig(host, port, user, pass, db, maxConns, minConns, time.Duration(maxLifetimeSec)*time.Second)

	log.Printf("[postgresql:user] config loaded: host=%s port=%s user=%s db=%s maxConns=%d minConns=%d lifetime=%s\n",
		host, port, user, db, maxConns, minConns, cfg.MaxConnLifetime)

	return cfg
}

func getEnv(key string) string {
	v := os.Getenv(key)

	if v == "" {
		log.Panicf("[postgresql:user][FATAL] env %s not set", key)
	}

	return v
}

func parseEnvInt32(key string) (int32, error) {
	valueStr := getEnv(key)

	value, err := strconv.ParseInt(valueStr, 10, 32)

	if err != nil {

		log.Panicf("[postgresql:user][FATAL] Invalid value for %s: %v", key, err)
	}

	return int32(value), nil
}
