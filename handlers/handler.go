package handlers

import (
	"ascii-art-web/art"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handleNotFound(w)
		return
	}
	if r.Method != http.MethodGet {
		handleMethodNotAllowed(w, r)
		return
	}
	loadAsciiArtForm(w)
}

func HandleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loadAsciiArtForm(w)
		return
	}

	if r.Method != http.MethodPost {
		handleMethodNotAllowed(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		handleBadRequest(w)
		return
	}
	text := r.FormValue("text")
	if text == "" {
		handleBadRequest(w)
		return
	}

	banner := r.FormValue("banner")
	if banner == "" {
		banner = "standard"
	}

	font, err := art.LoadBanner(banner)
	if err != nil {
		fmt.Println("Failed to load banner", err)
		handleInternalError(w)
		return
	}

	result, err := art.RenderText(text, font)
	if err != nil {
		fmt.Println(err)
		handleInternalError(w)
		return
	}

	tmpl := template.Must(template.ParseFiles(indexTemplatePath))
	err = tmpl.Execute(w, struct {
		Result         string
		SelectedBanner string
		OldInput       string
	}{
		Result:         result,
		SelectedBanner: banner,
		OldInput:       text})
	if err != nil {
		log.Println(err)
	}
}

func HandleError(w http.ResponseWriter, code int, msg string) {
	tmpl := template.Must(template.ParseFiles(errorTemplatePath))
	w.WriteHeader(code)
	err := tmpl.Execute(w, struct {
		ErrorCode    int
		ErrorMessage string
	}{
		code, msg,
	})
	if err != nil {
		log.Println(err)
	}
}

func loadAsciiArtForm(w http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles(indexTemplatePath))
	tmpl.Execute(w, nil)
}

func handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	responseText := fmt.Sprintf("Method %s is not allowed", r.Method)
	HandleError(w, http.StatusMethodNotAllowed, responseText)
}

func handleBadRequest(w http.ResponseWriter) {
	responseText := "Bad request"
	HandleError(w, http.StatusBadRequest, responseText)
}

func handleNotFound(w http.ResponseWriter) {
	responseText := "Not found"
	HandleError(w, http.StatusNotFound, responseText)
}

func handleInternalError(w http.ResponseWriter) {
	responseText := "Internal server error"
	HandleError(w, http.StatusInternalServerError, responseText)
}

const indexTemplatePath = "web/templates/index.html"
const errorTemplatePath = "web/templates/error.html"
