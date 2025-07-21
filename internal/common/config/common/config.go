package common

import (
	"github.com/1URose/marketplace/internal/common/settings"
	"log"
	"strings"
	"time"
)

type Config struct {
	GinAddress   string
	BearerPrefix string

	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewConfig(ginAddress, bearerPrefix, secretKey string, accessTTL, refreshTTL time.Duration) *Config {
	return &Config{
		GinAddress:   ginAddress,
		BearerPrefix: bearerPrefix,
		JWTSecret:    secretKey,
		AccessTTL:    accessTTL,
		RefreshTTL:   refreshTTL,
	}
}

func LoadCommonConfigFromEnv() *Config {
	const (
		envGinAddr      = "GIN_ADDRESS"
		envBearerPrefix = "AUTH_BEARER_PREFIX"
		envSecret       = "SECRET_KEY"
		envAccessTTL    = "ACCESS_TTL_MINUTES"  // в минутах
		envRefreshTTL   = "REFRESH_TTL_MINUTES" // в минутах
	)

	addr := settings.GetEnvSrt(envGinAddr)
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	log.Printf("[server:config] loaded GIN_ADDRESS=%s", addr)

	bearer := settings.GetEnvSrt(envBearerPrefix)
	if !strings.HasSuffix(bearer, " ") {
		bearer = bearer + " "
	}
	log.Printf("[server:config] loaded AUTH_BEARER_PREFIX=%q", bearer)

	secretKey := settings.GetEnvSrt(envSecret)
	log.Println("[server:config] loaded SECRET_KEY from env")

	accessSec, err := settings.GetEnvInt(envAccessTTL)
	if err != nil {
		log.Panicf("[server:config][FATAL] invalid %s: %v", envAccessTTL, err)
	}
	refreshSec, err := settings.GetEnvInt(envRefreshTTL)
	if err != nil {
		log.Panicf("[server:config][FATAL] invalid %s: %v", envRefreshTTL, err)
	}
	accessTTL := time.Duration(accessSec) * time.Minute
	refreshTTL := time.Duration(refreshSec) * time.Minute
	log.Printf("[server:config] token TTLs: AccessTTL=%s, RefreshTTL=%s", accessTTL, refreshTTL)

	return NewConfig(addr, bearer, secretKey, accessTTL, refreshTTL)
}
