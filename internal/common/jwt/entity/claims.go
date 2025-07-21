package entity

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}
