package main

import (
	"ascii-art-web/controller"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", controller.HandleHome)
	http.HandleFunc("/ascii-art", controller.HandleAsciiArt)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
	log.Println("Server stopped")

}
