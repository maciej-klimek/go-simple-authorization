// handlers/register.go
package handlers

import (
	"net/http"
	"simpleAuth/services"
	"simpleAuth/utils"
)

func register(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Register Handler called")
	if req.Method == http.MethodGet {
		Log.Debug("Serving registration page")
		http.ServeFile(wrt, req, "./static/html/register.html")
		return
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	Log.Debug("Parsed Form Data:")
	Log.Debugf("    - Email: %s", email)
	Log.Debugf("    - Password: %s", password)

	if !utils.CheckValidEmail(email) {
		Log.Warn("Invalid email provided:", email)
		http.Error(wrt, "Invalid email", http.StatusNotAcceptable)
		return
	}

	_, err := services.LoadUserData(email)
	if err == nil {
		Log.Warn("User already exists:", email)
		http.Error(wrt, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, _ := utils.HashPassword(password)
	user := services.LoginData{
		PasswordHash: hashedPassword,
	}
	err = services.SaveUserData(email, user)
	if err != nil {
		Log.Error("Failed to save user data:", err)
		http.Error(wrt, "Internal server error", http.StatusInternalServerError)
		return
	}

	Log.Info("User registered successfully")
	http.Redirect(wrt, req, "/login", http.StatusFound)
}
