package main

import (
	"log"
	"net/http"
	"simpleAuth/handlers"
	"simpleAuth/services"
	"simpleAuth/utils"
)

func main() {
	err := services.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	utils.Logger.Info("================= SERVER IS UP =================")
	handlers.Routes()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":8080", nil)
}
