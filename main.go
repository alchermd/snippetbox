package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	serverPort = ":4000"
)

// Define a home handler function which writes a welcome message to the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// Disable the catch-all behavior.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "hello from snippetbox")
}

// Handles showing of a snippet.
func showSnippet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Display a specific snippet")
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

func main() {
	// Create a new servemux and apply handler mappings.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Start a server on the given port and logging any potential error.
	log.Printf("Starting server on %s", serverPort)
	err := http.ListenAndServe(serverPort, mux)
	log.Fatal(err)
}
