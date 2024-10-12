package main

import (
	"net/http"
	"time"
)

func register(wrt http.ResponseWriter, req *http.Request) {
	log.Info("> Register Handler called")
	if req.Method != http.MethodPost {
		http.Error(wrt, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	if _, ok := Users[username]; ok {
		http.Error(wrt, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, _ := hashPassword(password)
	Users[username] = LoginData{
		PasswordHash: hashedPassword,
	}
	saveUserData()

	wrt.Write([]byte("User registered successfully"))
}

func login(wrt http.ResponseWriter, req *http.Request) {
	log.Info("Login Handler called")
	if req.Method != http.MethodPost {
		http.Error(wrt, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	user, ok := Users[username]

	if !ok || !checkPasswordHash(password, user.PasswordHash) {
		http.Error(wrt, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken := generateToken(32)

	http.SetCookie(wrt, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})
	user.SessionToken = sessionToken

	csrfToken := generateToken(32)

	http.SetCookie(wrt, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: false,
	})
	user.CSRFToken = csrfToken

	Users[username] = user

	wrt.Write([]byte("Login successful!"))
}

func content(wrt http.ResponseWriter, req *http.Request) {
	log.Info("> Chat Handler called")
	if req.Method != http.MethodPost {
		http.Error(wrt, "Invalid request method.", http.StatusMethodNotAllowed)
		return
	}

	if err := Authorize(req); err != nil {
		http.Error(wrt, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username := req.FormValue("username")
	wrt.Write([]byte("CSRF Token validation successful. Welcome, " + username))
}

func logout(wrt http.ResponseWriter, req *http.Request) {
	if err := Authorize(req); err != nil {
		http.Error(wrt, "Unauthorized", http.StatusUnauthorized)
		return
	}

	http.SetCookie(wrt, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(wrt, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	username := req.FormValue("username")
	user := Users[username]
	user.SessionToken = ""
	user.CSRFToken = ""
	Users[username] = user

	wrt.Write([]byte("Logged out successfully"))
}
