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
		err := "the resource at " + r.URL.Path + "that you are trying to access is not found"
		handleNotFound(w, err)
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
		handleBadRequest(w, "malformed request")
		return
	}
	text := r.FormValue("text")
	if text == "" {
		handleBadRequest(w, "empty input is not valid.")
		return
	}

	if len(text) > 1024 {
		handleBadRequest(w, "input is too long")
		return
	}

	banner := r.FormValue("banner")
	if banner == "" {
		banner = "standard"
	}

	font, err := art.LoadBanner(banner)
	if err != nil {
		log.Println("Failed to load banner", err)
		handleInternalError(w)
		return
	}

	result, err := art.RenderText(text, font)
	if err != nil {
		handleBadRequest(w, err.Error())
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
	str := "Method %s is not allowed"
	responseText := fmt.Sprintf(str, r.Method)
	HandleError(w, http.StatusMethodNotAllowed, responseText)
}

func handleBadRequest(w http.ResponseWriter, msg string) {
	responseText := "Bad request"
	if msg != "" {
		responseText += ": " + msg
	}
	HandleError(w, http.StatusBadRequest, responseText)
}

func handleNotFound(w http.ResponseWriter, msg string) {
	responseText := "Not found"
	if msg != "" {
		responseText += ": " + msg
	}
	HandleError(w, http.StatusNotFound, responseText)
}

func handleInternalError(w http.ResponseWriter) {
	responseText := "Internal server error"
	HandleError(w, http.StatusInternalServerError, responseText)
}

const indexTemplatePath = "web/templates/index.html"
const errorTemplatePath = "web/templates/error.html"
