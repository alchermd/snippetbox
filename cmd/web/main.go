package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/alchermd/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// This struct acts as a container for the shared dependencies of the application.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {
	// Load CLI options
	serverPort := flag.String("addr", ":4000", "Port that the server runs on")
	dsn := flag.String("dsn", "web:p@ssword!@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Setup custom loggers.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Setup dependency injection via struct initialization.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	// Initialize a server struct to use the custom error logger.
	srv := &http.Server{
		Addr:     *serverPort,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Start a server on the given port and logging any potential error.
	infoLog.Printf("Starting server on %s", srv.Addr)
	errorLog.Fatal(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
