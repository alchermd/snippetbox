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
	fmt.Fprintf(w, "hello from snippetbox")
}

func main() {
	// Create a new servemux and apply handler mappings.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Start a server on the given port and logging any potential error.
	log.Printf("Starting server on %s", serverPort)
	err := http.ListenAndServe(serverPort, mux)
	log.Fatal(err)
}
