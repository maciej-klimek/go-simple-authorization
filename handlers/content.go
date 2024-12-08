package handlers

import (
	"io"
	"net/http"
	"os"
	"simpleAuth/services"
)

func content(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Content Handler called")

	// If it's a GET request, only check the session token
	if req.Method == http.MethodGet {
		Log.Debug("Serving content page")
		sessionCookie, err := req.Cookie("session_token")
		if err != nil || sessionCookie.Value == "" {
			Log.Warn("No valid session token cookie found. Redirecting to login page.")
			http.Redirect(wrt, req, "/login", http.StatusFound)
			return
		}

		Log.Info("Valid session token cookie found. Serving content.")
		// Serve the content page if session token is valid
		http.ServeFile(wrt, req, "./static/html/content.html")
		return
	}

	// Check if it's a POST request (file upload)
	if req.Method == http.MethodPost {
		Log.Debug("Handling POST request for file upload")

		// Authorize request by checking session and CSRF token
		if err := services.Authorize(req); err != nil {
			Log.Warn("Authorization failed:", err)
			http.Redirect(wrt, req, "/login", http.StatusFound)
			return
		}
		// Handle file upload
		err := req.ParseMultipartForm(10 << 20) // Limit file size to 10MB
		if err != nil {
			Log.Error("Error parsing multipart form:", err)
			http.Error(wrt, "Error processing file", http.StatusInternalServerError)
			return
		}

		// Process file (same logic as before)
		file, handler, err := req.FormFile("file")
		if err != nil {
			Log.Error("Error retrieving file:", err)
			http.Error(wrt, "Error retrieving file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

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

		Log.Infof("File %s uploaded successfully", handler.Filename)
		wrt.Write([]byte("File uploaded successfully!"))
		return
	}

	// Handle other HTTP methods if needed
}
