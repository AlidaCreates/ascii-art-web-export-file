package controller

import (
	"net/http"
	"text/template"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}
func handleAsciiArt(w http.ResponseWriter, r *http.Request) {
	// Handle the ASCII art request
	w.Write([]byte("Here is some ASCII art!"))
}
