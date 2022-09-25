package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Setup middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Create a new servemux and apply handler mappings.
	mux := mux.NewRouter()
	mux.HandleFunc("/", app.home).Methods("GET")
	mux.HandleFunc("/snippet/create", app.createSnippet).Methods("POST")
	mux.HandleFunc("/snippet/create", app.createSnippetForm).Methods("GET")
	mux.HandleFunc("/snippet/{id}", app.showSnippet)
	mux.HandleFunc("/user/login", app.loginUserForm).Methods("GET")
	mux.HandleFunc("/user/login", app.loginUser).Methods("POST")
	mux.HandleFunc("/user/signup", app.signupUserForm).Methods("GET")
	mux.HandleFunc("/user/signup", app.signupUser).Methods("POST")
	mux.HandleFunc("/user/logout", app.logoutUser).Methods("POST")

	// Serve static assets.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
