package handlers

import (
	"net/http"
	"simpleAuth/utils"
)

var Log = utils.Logger

func Routes() {
	http.HandleFunc("/", content)
	http.HandleFunc("/view", viewFile)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
}
