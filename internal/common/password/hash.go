package password

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	log.Println("[password] hashing password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {

		log.Printf("[password][ERROR] failed to hash password: %v", err)

		return "", err
	}

	log.Println("[password] password hashed successfully")

	return string(hash), nil
}
