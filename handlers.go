package main

import (
	"html/template"
	"net/http"
	"time"
)

func register(wrt http.ResponseWriter, req *http.Request) {
	log.Info("> Register Handler called")
	if req.Method == http.MethodGet {
		// Serve the registration page
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(wrt, "Error loading page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(wrt, nil)
		return
	}

	email := req.FormValue("email")

	if !checkValidEmail(email) {
		http.Error(wrt, "Invalid email", http.StatusNotAcceptable)
		return
	}

	password := req.FormValue("password")

	if _, ok := Users[email]; ok {
		http.Error(wrt, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, _ := hashPassword(password)
	Users[email] = LoginData{
		PasswordHash: hashedPassword,
	}
	saveUserData()

	wrt.Write([]byte("User registered successfully"))
}

func login(wrt http.ResponseWriter, req *http.Request) {
	log.Info("Login Handler called")
	if req.Method == http.MethodGet {
		// Serve the login page
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(wrt, "Error loading page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(wrt, nil)
		return
	}

	email := req.FormValue("email")
	if !checkValidEmail(email) {
		http.Error(wrt, "Invalid email", http.StatusNotAcceptable)
		return
	}

	password := req.FormValue("password")
	user, ok := Users[email]
	log.Debug(user)

	if !ok || !checkPasswordHash(password, user.PasswordHash) {
		http.Error(wrt, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	sessionToken := generateToken(32)

	http.SetCookie(wrt, &http.Cookie{
		Name:     "email",
		Value:    email,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})
	user.SessionToken = sessionToken
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

	Users[email] = user

	saveUserData()

	wrt.Write([]byte("Login successful!"))
}

func content(wrt http.ResponseWriter, req *http.Request) {
	log.Info("> Chat Handler called")
	if req.Method == http.MethodGet {
		// Serve the login page
		tmpl, err := template.ParseFiles("templates/content.html")
		if err != nil {
			http.Error(wrt, "Error loading page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(wrt, nil)
		return
	}

	if err := Authorize(req); err != nil {
		http.Error(wrt, "Unauthorized", http.StatusUnauthorized)
		return
	}

	email, err := req.Cookie("email")
	if err == nil {
		wrt.Write([]byte("CSRF Token validation successful. Welcome, " + email.Value))
	}
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

	email, err := req.Cookie("email")
	if err == nil {
		user := Users[email.Value]
		user.SessionToken = ""
		user.CSRFToken = ""
		Users[email.Value] = user
		saveUserData()
	}

	wrt.Write([]byte("Logged out successfully"))
}
