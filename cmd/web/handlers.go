package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a welcome message to the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// Disable the catch-all behavior.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/footer.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err = ts.Execute(w, nil); err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Handles showing of a snippet.
func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Showing snippet with ID=%d.", id)
}

// Handles creation of a snippet.
func createSnippet(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests.
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "Create a new snippet")
}
