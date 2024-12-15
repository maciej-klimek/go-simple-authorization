package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"simpleAuth/services"
)

func content(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Content Handler called")

	if req.Method == http.MethodGet {
		Log.Debug("Serving content page")
		sessionCookie, err := req.Cookie("session_token")
		if err != nil || sessionCookie.Value == "" {
			Log.Warn("No valid session token cookie found. Redirecting to login page.")
			http.Redirect(wrt, req, "/login", http.StatusFound)
			return
		}

		Log.Info("Valid session token cookie found. Serving content.")
		http.ServeFile(wrt, req, "./static/html/content.html")
		return
	}

	if req.Method == http.MethodPost {
		Log.Debug("Handling POST request for file upload")

		if err := services.Authorize(req); err != nil {
			Log.Warn("Authorization failed:", err)
			http.Redirect(wrt, req, "/login", http.StatusFound)
			return
		}
		err := req.ParseMultipartForm(10 << 20)
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

		emailCookie, err := req.Cookie("email")
		if err != nil {
			Log.Error("Error retrieving email cookie:", err)
			http.Error(wrt, "Error retrieving email", http.StatusInternalServerError)
			return
		}
		email := emailCookie.Value

		userDir := filepath.Join("./uploads", email)
		if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
			Log.Error("Error creating user directory:", err)
			http.Error(wrt, "Error creating user directory", http.StatusInternalServerError)
			return
		}

		dst := filepath.Join(userDir, handler.Filename)
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
}
