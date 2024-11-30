package main

import (
	"html/template"
	"net/http"
	"time"
)

func index(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Index Handler called")

	sessionCookie, err := req.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		Log.Warn("No valid session token cookie found. Redirecting to login page.")
		http.Redirect(wrt, req, "/login", http.StatusFound)
		return
	}

	Log.Info("Valid session token cookie found. Redirecting to content page.")
	http.Redirect(wrt, req, "/content", http.StatusFound)
}

func register(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Register Handler called")
	if req.Method == http.MethodGet {
		Log.Debug("Serving registration page")
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			Log.Error("Error loading template:", err)
			http.Error(wrt, "Error loading page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(wrt, nil)
		return
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	Log.Debug("Parsed Form Data:")
	Log.Debugf("    - Email: %s", email)
	Log.Debugf("    - Password: %s", password)

	if !checkValidEmail(email) {
		Log.Warn("Invalid email provided:", email)
		http.Error(wrt, "Invalid email", http.StatusNotAcceptable)
		return
	}

	if _, ok := Users[email]; ok {
		Log.Warn("User already exists:", email)
		http.Error(wrt, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, _ := hashPassword(password)
	Users[email] = LoginData{
		PasswordHash: hashedPassword,
	}
	saveUserData()

	Log.Info("User registered successfully")
	wrt.Write([]byte("User registered successfully"))
}

func login(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Login Handler called")
	if req.Method == http.MethodGet {
		Log.Debug("Serving login page")
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			Log.Error("Error loading template:", err)
			http.Error(wrt, "Error loading page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(wrt, nil)
		return
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	Log.Debug("Parsed Form Data:")
	Log.Debugf("    - Email: %s", email)
	Log.Debugf("    - Password: %s", password)

	user, ok := Users[email]
	//Log.Debug("Retrieved User Data:", user)

	if !ok || !checkPasswordHash(password, user.PasswordHash) {
		Log.Warn("Invalid email or password")
		http.Error(wrt, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	Log.Debug("Generated Tokens:")
	Log.Debugf("    - Session Token: %s", sessionToken)
	Log.Debugf("    - CSRF Token: %s", csrfToken)

	http.SetCookie(wrt, &http.Cookie{
		Name:     "email",
		Value:    email,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(wrt, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})
	http.SetCookie(wrt, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: false,
	})
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	Users[email] = user

	saveUserData()

	Log.Info("Login successful")
	wrt.Write([]byte("Login successful!"))
}

func content(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Content Handler called")
	if req.Method == http.MethodGet {
		Log.Debug("Serving content page")
		tmpl, err := template.ParseFiles("templates/content.html")
		if err != nil {
			Log.Error("Error loading template:", err)
			http.Error(wrt, "Error loading page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(wrt, nil)
		return
	}

	if err := Authorize(req); err != nil {
		Log.Warn("Unauthorized request")
		http.Error(wrt, "Unauthorized", http.StatusUnauthorized)
		return
	}

	email, err := req.Cookie("email")
	if err == nil {
		Log.Debugf("CSRF Token validated for user: %s", email.Value)
		wrt.Write([]byte("CSRF Token validation successful. Welcome, " + email.Value))
	}
}

func logout(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Logout Handler called")

	// Check if the request is authorized
	if err := Authorize(req); err != nil {
		Log.Warn("Unauthorized request during logout")
		http.Error(wrt, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Clear session and CSRF cookies
	Log.Debug("Clearing session and CSRF cookies")

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

	http.SetCookie(wrt, &http.Cookie{
		Name:     "email",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
	})

	// Log user out from the session
	email, err := req.Cookie("email")
	if err == nil {
		Log.Debugf("Logging out user: %s", email.Value)
		user := Users[email.Value]
		user.SessionToken = ""
		user.CSRFToken = ""
		Users[email.Value] = user
		saveUserData() // Save updated user data
	} else {
		Log.Warn("No email cookie found during logout")
	}

	Log.Info("Logged out successfully")
	wrt.Write([]byte("Logged out successfully"))
}
