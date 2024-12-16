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

	if _, ok := Users[email]; ok {
		Log.Warn("User already exists:", email)
		http.Error(wrt, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, _ := utils.HashPassword(password)
	Users[email] = services.LoginData{
		PasswordHash: hashedPassword,
	}
	services.SaveUserData()

	Log.Info("User registered successfully")
	wrt.Write([]byte("User registered successfully"))
}
