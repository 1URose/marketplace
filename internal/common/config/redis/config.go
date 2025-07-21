package redis

import (
	"github.com/1URose/marketplace/internal/common/settings"
	"log"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       int
}

func NewConfig(host, port, user, password string, DB int) *Config {
	log.Printf("[redis:config] loading: host=%s port=%s user=%s db=%d", host, port, user, DB)
	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DB:       DB,
	}
}

func LoadRedisConfigFromEnv() *Config {
	log.Println("[redis:config] reading Redis config from env")

	const (
		envHost     = "REDIS_HOST"
		envPort     = "REDIS_PORT"
		envUser     = "REDIS_USER"
		envPassword = "REDIS_PASSWORD"
		envDB       = "REDIS_DB"
	)

	host := settings.GetEnvSrt(envHost)
	port := settings.GetEnvSrt(envPort)
	user := settings.GetEnvSrt(envUser)
	pass := settings.GetEnvSrt(envPassword)

	db, err := settings.GetEnvInt(envDB)
	if err != nil {
		log.Panicf("[redis:config][FATAL] invalid %s: %v", envDB, err)
	}

	log.Printf("[redis:config] loaded: host=%s port=%s user=%s db=%d", host, port, user, db)
	return NewConfig(host, port, user, pass, db)
}
