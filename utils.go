package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Error comparing password and hash:", err)
		return false
	}
	return true
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		Log.Fatalf("Failed to generate session token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)

}

func checkValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
