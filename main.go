package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("================== Server is up ==================\n ")
	loadUserData()
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/content", content)
	http.ListenAndServe(":8080", nil)

}
