package jwt

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

const (
	AccessTTL  = 1 * time.Hour
	RefreshTTL = 24 * time.Hour
)

func GenerateAccessToken(email string) (string, error) {

	log.Printf("[jwt] GenerateAccessToken called for email=%q", email)

	secret := GetSecretKeyFromEnv()

	claims := Claims{
		Email:     email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))

	if err != nil {

		log.Printf("[jwt][ERROR] GenerateAccessToken signing failed: %v", err)

		return "", err
	}

	log.Println("[jwt] GenerateAccessToken successful")

	return signed, nil
}

func GenerateRefreshToken(email string) (string, error) {

	log.Printf("[jwt] GenerateRefreshToken called for email=%q", email)

	secret := GetSecretKeyFromEnv()

	claims := Claims{
		Email:     email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))

	if err != nil {

		log.Printf("[jwt][ERROR] GenerateRefreshToken signing failed: %v", err)

		return "", err
	}

	log.Println("[jwt] GenerateRefreshToken successful")

	return signed, nil
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := parseToken(tokenString, GetSecretKeyFromEnv())
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "access" {
		return nil, fmt.Errorf("expected access token, got %q", claims.TokenType)
	}
	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := parseToken(tokenString, GetSecretKeyFromEnv())
	if err != nil {
		return nil, err
	}
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("expected refresh token, got %q", claims.TokenType)
	}
	return claims, nil
}

func parseToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}
	return claims, nil
}

func GetSecretKeyFromEnv() string {

	key := getEnv("SECRET_KEY")

	log.Println("[jwt] SECRET_KEY loaded from env")

	return key
}

func getEnv(key string) string {

	value := os.Getenv(key)

	if value == "" {

		log.Panicf("[jwt][FATAL] environment variable %s is not set", key)
	}

	return value
}
