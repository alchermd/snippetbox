package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alchermd/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

// This struct acts as a container for the shared dependencies of the application.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	session       *sessions.CookieStore
}

func main() {
	// Load CLI options
	serverPort := flag.String("addr", ":4000", "Port that the server runs on")
	dsn := flag.String("dsn", "web:p@ssword!@/snippetbox?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	// Setup custom loggers.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Setup sessions
	session := sessions.NewCookieStore([]byte(*secret))
	session.Options.MaxAge = int(time.Hour) * 12
	session.Options.Secure = true

	// Setup dependency injection via struct initialization.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
	}

	// Initialize a server struct to use the custom error logger.
	srv := &http.Server{
		Addr:     *serverPort,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Start a server on the given port and logging any potential error.
	infoLog.Printf("Starting server on %s", srv.Addr)
	errorLog.Fatal(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
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
