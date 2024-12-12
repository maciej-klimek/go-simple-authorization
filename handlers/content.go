package handlers

import (
	"net/http"
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
		Log.Debug("Handling POST request")

		if err := services.Authorize(req); err != nil {
			Log.Warn("Authorization failed:", err)
			http.Redirect(wrt, req, "/login", http.StatusFound)
			return
		}

		Log.Infof("POST request handled successfully")
		wrt.Write([]byte("POST request handled successfully!"))
		return
	}
}
