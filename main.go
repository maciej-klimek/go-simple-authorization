package main

import (
	"net/http"
)

func main() {

	ConfigureLogger("info")
	log.Info("================== Server is up ==================")

	loadUserData()

	http.HandleFunc("/", index)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/content", content)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // Serve static files (CSS, JS)
	http.ListenAndServe(":8080", nil)

}

func index(wrt http.ResponseWriter, req *http.Request) {
	http.Redirect(wrt, req, "/login", http.StatusFound)
}
