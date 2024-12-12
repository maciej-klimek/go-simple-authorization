package handlers

import (
	"net/http"
	"simpleAuth/services"
	"simpleAuth/utils"
	"time"
)

func login(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Login Handler called")

	if req.Method == http.MethodGet {
		Log.Debug("Serving login page")
		http.ServeFile(wrt, req, "./static/html/login.html")
		return
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	Log.Debug("Parsed Form Data:")
	Log.Debugf("    - Email: %s", email)
	Log.Debugf("    - Password: %s", password)

	user, err := services.LoadUserData(email)
	if err != nil || !utils.CheckPasswordHash(password, user.PasswordHash) {
		Log.Warn("Invalid email or password")
		http.Error(wrt, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	sessionToken := utils.GenerateToken(32)
	csrfToken := utils.GenerateToken(32)

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
	err = services.SaveUserData(email, user)
	if err != nil {
		Log.Error("Failed to save user data:", err)
		http.Error(wrt, "Internal server error", http.StatusInternalServerError)
		return
	}

	Log.Info("Login successful")
	http.Redirect(wrt, req, "/content", http.StatusFound)
}
