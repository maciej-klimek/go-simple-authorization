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

	fileName := req.URL.Query().Get("file")
	if fileName == "" {
		http.Error(wrt, "File not specified", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("/shared-data", email, fileName)

	fileInfo, err := os.Stat(filePath)
	if err != nil || fileInfo.IsDir() {
		http.Error(wrt, "File not found", http.StatusNotFound)
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		// Fallback to text/plain for .txt files or unknown types
		mimeType = "application/octet-stream"
		if filepath.Ext(filePath) == ".txt" {
			mimeType = "text/plain"
		}
	}

	wrt.Header().Set("Content-Type", mimeType)
	if mimeType == "application/octet-stream" {
		wrt.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	} else {
		wrt.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(filePath))
	}

	file, err := os.Open(filePath)
	if err != nil {
		Log.Error("Error opening file:", err)
		http.Error(wrt, "Error accessing file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	http.ServeContent(wrt, req, fileInfo.Name(), fileInfo.ModTime(), file)
}
