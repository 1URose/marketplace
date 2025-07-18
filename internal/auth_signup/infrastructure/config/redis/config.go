package redis

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       int
}

func NewConfig(host, port, user, password string, DB int) *Config {
	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DB:       DB,
	}
}

func ReadConfigFromEnv() *Config {
	log.Println("[redis:config] reading Redis configuration from env")

	host := getEnv("REDIS_HOST")
	port := getEnv("REDIS_PORT")
	user := getEnv("REDIS_USER")
	password := getEnv("REDIS_PASSWORD")

	db, err := parseEnvInt("REDIS_DB")
	if err != nil {
		log.Printf("[redis:config][ERROR] failed to parse REDIS_DB: %v", err)
		panic(err)
	}

	log.Printf("[redis:config] config loaded: host=%s port=%s user=%s db=%d",
		host, port, user, db)
	return NewConfig(host, port, user, password, db)
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("[redis:config][FATAL] environment variable %s is not set", key)
	}
	return value
}

func parseEnvInt(key string) (int, error) {
	s := getEnv(key)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid value for %s: %w", key, err)
	}
	return v, nil
}
