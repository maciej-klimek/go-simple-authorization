package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
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

	// Check if the session token is present and valid
	sessionCookie, err := req.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		Log.Warn("No valid session token. Redirecting to login page.")
		http.Redirect(wrt, req, "/login", http.StatusFound)
		return
	}
	Log.Debugf("Session token found: %s", sessionCookie.Value)

	// Check for the email cookie
	emailCookie, err := req.Cookie("email")
	if err != nil || emailCookie.Value == "" {
		Log.Warn("No email cookie. Redirecting to login page.")
		http.Redirect(wrt, req, "/login", http.StatusFound)
		return
	}
	Log.Debugf("Email cookie found: %s", emailCookie.Value)

	// Retrieve user from the Users map using the email
	user, ok := Users[emailCookie.Value]
	if !ok {
		Log.Warn("User not found in Users map for email:", emailCookie.Value)
		http.Redirect(wrt, req, "/login", http.StatusFound)
		return
	}
	Log.Debugf("User session token: %s", user.SessionToken)

	// Validate the session token
	if user.SessionToken != sessionCookie.Value {
		Log.Warn("Session token mismatch for user:", emailCookie.Value)
		http.Redirect(wrt, req, "/login", http.StatusFound)
		return
	}
	Log.Debug("Session token validated successfully.")

	// // CSRF Token validation

	// for key, values := range req.Header {
	// 	for _, value := range values {
	// 		Log.Debugf("Header: %s: %s", key, value)
	// 	}
	// }

	// csrfTokenHeader := req.Header.Get("X-CSRF-Token")
	// csrfTokenCookie, err := req.Cookie("csrf_token")
	// if err != nil {
	// 	Log.Warn("Error retrieving CSRF token from cookies:", err)
	// 	http.Error(wrt, "CSRF token missing", http.StatusUnauthorized)
	// 	return
	// }
	// Log.Debugf("CSRF Token from header: %s", csrfTokenHeader)
	// Log.Debugf("CSRF Token from cookie: %s", csrfTokenCookie.Value)

	// // Check if CSRF tokens match
	// if csrfTokenHeader != csrfTokenCookie.Value {
	// 	Log.Warn("CSRF token mismatch: header:", csrfTokenHeader, "cookie:", csrfTokenCookie.Value)
	// 	http.Error(wrt, "Invalid CSRF token", http.StatusUnauthorized)
	// 	return
	// }
	// Log.Debug("CSRF token validated successfully.")

	// Handle GET request (serve content page)
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

	// Handle POST request for file upload
	if req.Method == http.MethodPost {
		Log.Debug("Handling file upload request.")

		// Parse the uploaded file
		err := req.ParseMultipartForm(10 << 20) // Limit file size to 10MB
		if err != nil {
			Log.Error("Error parsing multipart form:", err)
			http.Error(wrt, "Error processing file", http.StatusInternalServerError)
			return
		}

		file, handler, err := req.FormFile("file")
		if err != nil {
			Log.Error("Error retrieving file:", err)
			http.Error(wrt, "Error retrieving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Save the uploaded file to the server
		dst := "./uploads/" + handler.Filename
		out, err := os.Create(dst)
		if err != nil {
			Log.Error("Error saving file:", err)
			http.Error(wrt, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			Log.Error("Error copying file:", err)
			http.Error(wrt, "Error saving file", http.StatusInternalServerError)
			return
		}

		Log.Infof("File %s uploaded successfully by user %s", handler.Filename, emailCookie.Value)
		wrt.Write([]byte("File uploaded successfully!"))
	}
}

func logout(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Logout Handler called")

	if err := Authorize(req); err != nil {
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
		saveUserData() // Save updated user data
	} else {
		Log.Warn("No email cookie found during logout")
	}

	Log.Info("Logged out successfully")
	wrt.Write([]byte("Logged out successfully"))
}
