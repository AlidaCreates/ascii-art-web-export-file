package controller

import (
	"net/http"
	"text/template"
	"ascii-art-web/art"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}
func HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
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
	http.Error(w, "Failed to load banner", http.StatusInternalServerError)
	return
	}

	result, err := art.RenderText(text, font)
	if err != nil {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
	return
	}

	tmpl := template.Must(template.ParseFiles("templates/ascii.html"))
	tmpl.Execute(w, struct{ Result string }{Result: result})
}
