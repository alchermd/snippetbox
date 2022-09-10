package main

import (
	"log"
	"net/http"
)

var (
	serverPort = ":4000"
)

func main() {
	// Create a new servemux and apply handler mappings.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Serve static assets.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Start a server on the given port and logging any potential error.
	log.Printf("Starting server on %s", serverPort)
	log.Fatal(http.ListenAndServe(serverPort, mux))
}
