package controller

import (
	"ascii-art-web/art"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	loadAsciiArtForm(w)
}

func HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loadAsciiArtForm(w)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}
	banner := r.FormValue("banner")
	if banner == "" {
		http.Error(w, "Banner is required", http.StatusBadRequest)
		return
	}

	font, err := art.LoadBanner(banner)
	if err != nil {
		fmt.Println("Failed to load banner", err)
		http.Error(w, "Failed to load banner", http.StatusInternalServerError)
		return
	}

	result, err := art.RenderText(text, font)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	err = tmpl.Execute(w, struct{ Result string }{Result: result})
	if err != nil {
		log.Println(err)
	}
}

func loadAsciiArtForm(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	tmpl.Execute(w, nil)
}
