package main

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

var ErrAuth = errors.New("Unauthorized")

func Authorize(req *http.Request) error {
	email := req.FormValue("email")

	user, ok := Users[email]

	if !ok {
		log.Error("User not found")
		return ErrAuth
	}

	sessionToken, err := req.Cookie("session_token")
	if err != nil {
		log.WithFields(logrus.Fields{"error": err}).Info("Failed to get session token")
		return ErrAuth
	}

	log.WithFields(logrus.Fields{
		"request_token": sessionToken.Value,
		"db_token":      user.SessionToken,
	}).Debug("Session token check")

	if sessionToken.Value == "" || sessionToken.Value != user.SessionToken {
		log.Warn("Session token mismatch or empty")
		return ErrAuth
	}

	csrfToken := req.Header.Get("X-CSRF-Token")
	log.WithFields(logrus.Fields{
		"request_token": csrfToken,
		"db_token":      user.CSRFToken,
	}).Debug("CSRF token check")

	if csrfToken != user.CSRFToken {
		log.Warn("CSRF token mismatch")
		return ErrAuth
	}

	return nil
}
