package handlers

import (
	"net/http"
	"simpleAuth/services"
	"simpleAuth/utils"
)

var Log = utils.Logger
var Users = services.Users

func Routes() {
	http.HandleFunc("/", content)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
}
