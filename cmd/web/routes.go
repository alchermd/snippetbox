package main

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Setup middleware chain
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddelware := alice.New(app.authenticate)

	// Create a new servemux and apply handler mappings.
	mux := mux.NewRouter()
	mux.Handle("/", dynamicMiddelware.ThenFunc(app.home)).Methods("GET")
	mux.Handle("/snippet/create", dynamicMiddelware.Append(app.requireAuthentication).ThenFunc(app.createSnippet)).Methods("POST")
	mux.Handle("/snippet/create", dynamicMiddelware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm)).Methods("GET")
	mux.Handle("/snippet/{id}", dynamicMiddelware.ThenFunc(app.showSnippet)).Methods("GET")
	mux.Handle("/user/login", dynamicMiddelware.ThenFunc(app.loginUserForm)).Methods("GET")
	mux.Handle("/user/login", dynamicMiddelware.ThenFunc(app.loginUser)).Methods("POST")
	mux.Handle("/user/signup", dynamicMiddelware.ThenFunc(app.signupUserForm)).Methods("GET")
	mux.Handle("/user/signup", dynamicMiddelware.ThenFunc(app.signupUser)).Methods("POST")
	mux.Handle("/user/logout", dynamicMiddelware.ThenFunc(app.logoutUser)).Methods("POST")

	// Serve static assets.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return csrf.Protect([]byte("32-byte-long-auth-key"))(standardMiddleware.Then(mux))
}
