package handlers

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"simpleAuth/services"
)

func content(wrt http.ResponseWriter, req *http.Request) {
	Log.Info("Content Handler called")

	// GET request: Serve content page and list user files
	if req.Method == http.MethodGet {
		Log.Debug("Serving content page")
		sessionCookie, err := req.Cookie("session_token")
		if err != nil || sessionCookie.Value == "" {
			Log.Warn("No valid session token cookie found. Redirecting to login page.")
			http.Redirect(wrt, req, "/login", http.StatusFound)
			return
		}

		Log.Info("Valid session token cookie found. Serving content.")

		// Retrieve the logged-in user's email to list their files
		emailCookie, err := req.Cookie("email")
		if err != nil {
			Log.Error("Error retrieving email cookie:", err)
			http.Error(wrt, "Error retrieving user data", http.StatusInternalServerError)
			return
		}
		email := emailCookie.Value

		// Construct the user's directory path
		userDir := filepath.Join("/shared-data", email)

		// List all files in the user's folder
		files, err := os.ReadDir(userDir)
		if err != nil {
			Log.Error("Error reading user directory:", err)
			http.Error(wrt, "Error accessing user files", http.StatusInternalServerError)
			return
		}

		// Prepare the file list for injection into the HTML
		var fileList []string
		for _, file := range files {
			if !file.IsDir() {
				// Add each file name to the list
				fileList = append(fileList, file.Name())
			}
		}

		// Create a template and pass the file list to the template
		tmpl, err := template.ParseFiles("./static/html/content.html")
		if err != nil {
			Log.Error("Error parsing template:", err)
			http.Error(wrt, "Error rendering page", http.StatusInternalServerError)
			return
		}

		// Render the page with the file list
		err = tmpl.Execute(wrt, struct {
			Files []string
		}{Files: fileList})

		if err != nil {
			Log.Error("Error rendering template:", err)
			http.Error(wrt, "Error rendering page", http.StatusInternalServerError)
			return
		}

		return
	}

	// POST request: Handle file upload
	if req.Method == http.MethodPost {
		Log.Debug("Handling POST request for file upload")

		if err := services.Authorize(req); err != nil {
			Log.Warn("Authorization failed:", err)
			http.Redirect(wrt, req, "/login", http.StatusFound)
			return
		}

		emailCookie, err := req.Cookie("email")
		if err != nil {
			log.Fatal(err)
		}
		email := emailCookie.Value

		err = req.ParseMultipartForm(10 << 20)
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

		dst := filepath.Join("/shared-data/"+email, handler.Filename)
		err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

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
