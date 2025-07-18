package password

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CheckPasswordHash(password, hash string) bool {
	log.Println("[password] verifying password hash")

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {

		log.Printf("[password][ERROR] password verification failed: %v", err)

		return false
	}

	log.Println("[password] password verification succeeded")

	return true
}
