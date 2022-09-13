package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// This struct acts as a container for the shared dependencies of the application.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Load CLI options
	serverPort := flag.String("addr", ":4000", "Port that the server runs on")
	flag.Parse()

	// Setup custom loggers.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Setup dependency injection via struct initialization.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Create a new servemux and apply handler mappings.
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Serve static assets.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Initialize a server struct to use the custom error logger.
	srv := &http.Server{
		Addr:     *serverPort,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Start a server on the given port and logging any potential error.
	infoLog.Printf("Starting server on %s", srv.Addr)
	errorLog.Fatal(srv.ListenAndServe())
}
