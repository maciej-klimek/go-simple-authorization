package main

import (
	"errors"
	"net/http"
)

var AuthError = errors.New("Unauthorized")

func Authorize(req *http.Request) error {
	username := req.FormValue("username")
	user, ok := Users[username]
	if !ok {
		return AuthError
	}

	// Get session token from cookies
	sessionToken, err := req.Cookie("session_token")
	if err != nil || sessionToken.Value == "" || sessionToken.Value == user.SessionToken {
		return AuthError
	}

	// Get csrf token from the header !! (thats the real auth here)
	csrfToken := req.Header.Get("X-CSRFToken")
	if csrfToken != user.CSRFToken {
		return AuthError
	}

	return nil
}
