package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/alchermd/snippetbox/pkg/models"
	"github.com/gorilla/mux"
)

// Define a home handler function which writes a welcome message to the response body.
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: snippets,
	})
}

// Handles showing of a snippet.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: snippet,
	})
}

// Handles creation of a snippet.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank."
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum 100 characters)."
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank."
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be blank."
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = `This field is invalid (must be either "365", "7", or "1").`
	}

	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}
