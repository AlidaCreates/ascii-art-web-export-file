package main

import (
	"ascii-art-web/handlers"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handlers.HandleHome)
	http.HandleFunc("/ascii-art", handlers.HandleAsciiArt)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
	log.Println("Server stopped")

}
