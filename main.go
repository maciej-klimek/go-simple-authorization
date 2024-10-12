package main

import (
	"net/http"
)

func main() {

	ConfigureLogger("debug")
	log.Info("================== Server is up ==================")

	loadUserData()

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/content", content)
	http.ListenAndServe(":8080", nil)
}
