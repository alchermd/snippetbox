package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alchermd/snippetbox/pkg/models"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			session, _ := app.session.Get(r, "session-name")
			session.Values["flash"] = "You need to login first."

			if err := session.Save(r, w); err != nil {
				app.serverError(w, err)
				return
			}

			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.session.Get(r, "session-name")
		if err != nil {
			app.serverError(w, err)
			return
		}

		if session.Values["authenticatedUserID"] == nil {
			next.ServeHTTP(w, r)
			return
		}

		id, err := strconv.Atoi(fmt.Sprintf("%v", session.Values["authenticatedUserID"]))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.users.Get(id)
		if errors.Is(err, models.ErrNoRecord) {
			delete(session.Values, "authenticatedUserID")
			if err := session.Save(r, w); err != nil {
				app.serverError(w, err)
				return
			}

			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		if !user.Active {
			delete(session.Values, "authenticatedUserID")
			if err := session.Save(r, w); err != nil {
				app.serverError(w, err)
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
