package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Error comparing password and hash:", err)
		return false
	}
	return true
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		Logger.Fatalf("Failed to generate session token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)

}

func CheckValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
