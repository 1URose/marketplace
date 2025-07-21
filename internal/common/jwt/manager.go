package jwt

import (
	"fmt"
	"github.com/1URose/marketplace/internal/common/config/common"
	"github.com/1URose/marketplace/internal/common/jwt/entity"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Manager struct {
	secretKey  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewManager(cfg *common.Config) *Manager {
	log.Printf("[jwt] NewManager called: accessTTL=%s refreshTTL=%s", cfg.AccessTTL, cfg.RefreshTTL)

	return &Manager{
		secretKey:  []byte(cfg.JWTSecret),
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
	}
}

func (m *Manager) GenerateAccessToken(email string, UserId int) (string, error) {

	log.Printf("[jwt] GenerateAccessToken called for email=%q", email)

	claims := entity.Claims{
		Email:     email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(UserId),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(m.secretKey)

	if err != nil {

		log.Printf("[jwt][ERROR] GenerateAccessToken signing failed: %v", err)

		return "", err
	}

	log.Println("[jwt] GenerateAccessToken successful")

	return signed, nil
}

func (m *Manager) GenerateRefreshToken(email string, UserId int) (string, error) {

	log.Printf("[jwt] GenerateRefreshToken called for email=%q", email)

	claims := entity.Claims{
		Email:     email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(UserId),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(m.secretKey)

	if err != nil {

		log.Printf("[jwt][ERROR] GenerateRefreshToken signing failed: %v", err)

		return "", err
	}

	log.Println("[jwt] GenerateRefreshToken successful")

	return signed, nil
}

func (m *Manager) ValidateAccessToken(tokenString string) (*entity.Claims, error) {
	log.Println("[jwt] ValidateAccessToken called")
	claims, err := m.parseToken(tokenString)
	if err != nil {
		log.Printf("[jwt][ERROR] ValidateAccessToken parseToken failed: %v", err)
		return nil, err
	}
	if claims.TokenType != "access" {
		err := fmt.Errorf("expected access token, got %q", claims.TokenType)
		log.Printf("[jwt][ERROR] ValidateAccessToken wrong token type: %v", err)
		return nil, err
	}
	log.Printf("[jwt] ValidateAccessToken successful: subject=%s email=%s", claims.Subject, claims.Email)
	return claims, nil
}

func (m *Manager) ValidateRefreshToken(tokenString string) (*entity.Claims, error) {
	log.Println("[jwt] ValidateRefreshToken called")
	claims, err := m.parseToken(tokenString)
	if err != nil {
		log.Printf("[jwt][ERROR] ValidateRefreshToken parseToken failed: %v", err)
		return nil, err
	}
	if claims.TokenType != "refresh" {
		err := fmt.Errorf("expected refresh token, got %q", claims.TokenType)
		log.Printf("[jwt][ERROR] ValidateRefreshToken wrong token type: %v", err)
		return nil, err
	}
	log.Printf("[jwt] ValidateRefreshToken successful: subject=%s email=%s", claims.Subject, claims.Email)
	return claims, nil
}

func (m *Manager) parseToken(tokenString string) (*entity.Claims, error) {
	log.Printf("[jwt] parseToken called")
	token, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			err := fmt.Errorf("unexpected signing method %v", t.Header["alg"])
			log.Printf("[jwt][ERROR] parseToken unexpected signing method: %v", err)
			return nil, err
		}
		return m.secretKey, nil
	})
	if err != nil {
		log.Printf("[jwt][ERROR] parseToken parsing failed: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*entity.Claims)
	if !ok || !token.Valid {
		log.Printf("[jwt][ERROR] parseToken invalid token or claims assertion failed")
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		log.Printf("[jwt][ERROR] parseToken token expired at %s", claims.ExpiresAt.Time.Format(time.RFC3339))
		return nil, fmt.Errorf("token expired")
	}

	log.Printf(
		"[jwt] parseToken successful: tokenType=%s subject=%s email=%s issuedAt=%s expiresAt=%s",
		claims.TokenType,
		claims.Subject,
		claims.Email,
		claims.IssuedAt.Time.Format(time.RFC3339),
		claims.ExpiresAt.Time.Format(time.RFC3339),
	)
	return claims, nil
}
