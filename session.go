package main

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

var ErrAuth = errors.New("Unauthorized")

func Authorize(req *http.Request) error {
	email, err := req.Cookie("email")

	user, ok := Users[email.Value]

	if !ok {
		Log.Error("User not found")
		return ErrAuth
	}

	sessionToken, err := req.Cookie("session_token")
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Info("Failed to get session token")
		return ErrAuth
	}

	Log.WithFields(logrus.Fields{
		"request_token": sessionToken.Value,
		"db_token":      user.SessionToken,
	}).Debug("Session token check")

	if sessionToken.Value == "" || sessionToken.Value != user.SessionToken {
		Log.Warn("Session token mismatch or empty")
		return ErrAuth
	}

	csrfToken, err := req.Cookie("csrf_token")
	Log.WithFields(logrus.Fields{
		"request_token": csrfToken,
		"db_token":      user.CSRFToken,
	}).Debug("CSRF token check")

	if csrfToken.Value != user.CSRFToken {
		Log.Warn("CSRF token mismatch")
		return ErrAuth
	}

	return nil
}
