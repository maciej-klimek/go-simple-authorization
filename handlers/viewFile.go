package handlers

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"simpleAuth/services"
)

func viewFile(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("View Handler called")

	// Authorize the user
	err := services.Authorize(req)
	if err != nil {
		Log.Warn("Authorization failed:", err)
		http.Error(wrt, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	emailCookie, err := req.Cookie("email")
	if err != nil || emailCookie.Value == "" {
		Log.Error("Error retrieving email cookie or cookie is empty")
		http.Error(wrt, "Unauthorized access", http.StatusUnauthorized)
		return
	}
	email := emailCookie.Value

	// Get the file name from the query parameter
	fileName := req.URL.Query().Get("file")
	if fileName == "" {
		http.Error(wrt, "File not specified", http.StatusBadRequest)
		return
	}

	// Construct the full file path
	filePath := filepath.Join("/shared-data", email, fileName)

	// Check if the file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.IsDir() {
		http.Error(wrt, "File not found", http.StatusNotFound)
		return
	}

	// Determine the MIME type of the file
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		// Fallback to text/plain for .txt files or unknown types
		mimeType = "application/octet-stream"
		if filepath.Ext(filePath) == ".txt" {
			mimeType = "text/plain"
		}
	}

	// Set the appropriate headers
	wrt.Header().Set("Content-Type", mimeType)
	if mimeType == "application/octet-stream" {
		// For binary files or unknown types, force download
		wrt.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	} else {
		// Render files inline for supported types
		wrt.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	}

	// Serve the file
	file, err := os.Open(filePath)
	if err != nil {
		Log.Error("Error opening file:", err)
		http.Error(wrt, "Error accessing file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	http.ServeContent(wrt, req, fileInfo.Name(), fileInfo.ModTime(), file)
}
