package main

import (
	"net/http"
	"simpleAuth/handlers"
	"simpleAuth/services"
	"simpleAuth/utils"
)

func main() {

	utils.Logger.Info("================= SERVER IS UP =================")
	services.LoadUserData()
	handlers.Routes()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}
