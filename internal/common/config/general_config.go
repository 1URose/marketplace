package config

import (
	"github.com/1URose/marketplace/internal/common/config/ad_limits"
	"github.com/1URose/marketplace/internal/common/config/common"
	"github.com/1URose/marketplace/internal/common/config/postgresql"
	"github.com/1URose/marketplace/internal/common/config/redis"
	"github.com/joho/godotenv"
	"log"
)

type GeneralConfig struct {
	AdConfig       *ad_limits.AdConfig
	PostgresConfig *postgresql.Config
	RedisConfig    *redis.Config
	CommonConfig   *common.Config
}

func NewGeneralConfig() *GeneralConfig {
	log.Println("Creating GeneralConfig")
	return &GeneralConfig{
		AdConfig:       ad_limits.LoadAdConfigFromEnv(),
		PostgresConfig: postgresql.LoadPGConfigFromEnv(),
		RedisConfig:    redis.LoadRedisConfigFromEnv(),
		CommonConfig:   common.LoadCommonConfigFromEnv(),
	}
}

func LoadGeneralConfigFrom(path string) (*GeneralConfig, error) {
	log.Printf("Loading environment from %s", path)
	if err := godotenv.Load(path); err != nil {
		log.Printf("ERROR: failed to load env file %s: %v", path, err)
		return nil, err
	}
	log.Println("Environment variables loaded successfully")
	cfg := NewGeneralConfig()
	return cfg, nil
}
