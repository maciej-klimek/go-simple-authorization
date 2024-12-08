package services

import (
	"errors"
	"net/http"

	"simpleAuth/utils"

	"github.com/sirupsen/logrus"
)

var Log = utils.Logger

var ErrAuth = errors.New("Unauthorized")

func Authorize(req *http.Request) error {
	emailCookie, err := req.Cookie("email")
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Error("Failed to get email cookie")
		return ErrAuth
	}

	email := emailCookie.Value
	user, ok := Users[email]
	if !ok {
		Log.Error("User not found")
		return ErrAuth
	}

	sessionCookie, err := req.Cookie("session_token")
	if err != nil {
		Log.WithFields(logrus.Fields{"error": err}).Info("Failed to get session token cookie")
		return ErrAuth
	}

	sessionToken := sessionCookie.Value
	Log.WithFields(logrus.Fields{
		"request_token": sessionToken,
		"local_token":   user.SessionToken,
	}).Debug("Session token check")

	if sessionToken == "" || sessionToken != user.SessionToken {
		Log.Warn("Session token mismatch or empty")
		return ErrAuth
	}

	//debugg := req.Header
	//Log.Debug(debugg)
	csrfToken := req.Header.Get("X-CSRF-Token")
	Log.WithFields(logrus.Fields{
		"request_token": csrfToken,
		"local_token":   user.CSRFToken,
	}).Debug("CSRF token check")

	if csrfToken != user.CSRFToken {
		Log.Warn("CSRF token mismatch")
		return ErrAuth
	}

	return nil
}
