package main

import (
	"fmt"
	"net/http"
	"time"
)

func register(wrt http.ResponseWriter, req *http.Request) {
	fmt.Println("> Register Handler called")
	if req.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(wrt, "Invalid methon", err)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	if _, ok := Users[username]; ok {
		err := http.StatusConflict
		http.Error(wrt, "User already exists", err)
		return
	}

	hashedPassword, _ := hashPassword(password)
	Users[username] = LoginData{
		PasswordHash: hashedPassword,
	}
	saveUserData()

	fmt.Fprintln(wrt, "User registered successfully")
}

func login(wrt http.ResponseWriter, req *http.Request) {
	fmt.Println("> Login Handler called")
	if req.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(wrt, "Invalid methon", err)
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")

	user, ok := Users[username]

	if !ok || !checkPasswordHash(password, user.PasswordHash) {
		err := http.StatusUnauthorized
		http.Error(wrt, "Invalid username or password", err)
		return
	}

	// setting sessionToken
	sessionToken := generateToken(32)

	http.SetCookie(wrt, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})
	user.SessionToken = sessionToken

	// setting csrfToken
	csrfToken := generateToken(32)

	http.SetCookie(wrt, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: false,
	})
	user.CSRFToken = csrfToken

	Users[username] = user

	fmt.Fprintln(wrt, "Login successful!")
}

func content(wrt http.ResponseWriter, req *http.Request) {
	fmt.Println("> Chat Handler called")
	if req.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(wrt, "Invalid request method.", err)
		return
	}

	if err := Authorize(req); err != nil {
		err := http.StatusUnauthorized
		http.Error(wrt, "Unauthorized", err)
		return
	}

	username := req.FormValue("username")
	fmt.Fprintf(wrt, "CSRF Token validation successful. Welcome, %s", username)
}

func logout(wrt http.ResponseWriter, req *http.Request) {
	if err := Authorize(req); err != nil {
		err := http.StatusUnauthorized
		http.Error(wrt, "Unauthorized", err)
		return
	}

	// Clear cookie
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

	fmt.Fprintln(wrt, "Logged out successfuly")
}
