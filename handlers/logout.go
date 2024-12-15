package handlers

import (
	"net/http"
	"simpleAuth/services"
	"time"
)

func logout(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Logout Handler called")

	if err := services.Authorize(req); err != nil {
		Log.Warn("Unauthorized request during logout")
		http.Error(wrt, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	email, err := req.Cookie("email")
	if err == nil {
		Log.Debugf("Logging out user: %s", email.Value)
		user := Users[email.Value]
		user.SessionToken = ""
		user.CSRFToken = ""
		Users[email.Value] = user
		services.SaveUserData() // Save updated user data
	} else {
		Log.Warn("No email cookie found during logout")
	}

	Log.Info("Logged out successfully")
	wrt.Write([]byte("Logged out successfully"))
}
