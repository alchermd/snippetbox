package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/csrf"
)

// Writes a generic server error response while logging debugging information to stderr.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Writes a response, intended to be used for client errors.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Convenience wrapper for 404 errors
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Fetch the given template name from the cache and execute it with the given template data.
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, w, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, w http.ResponseWriter, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()

	session, _ := app.session.Get(r, "session-name")
	storedFlash, ok := session.Values["flash"]
	flashMessage := ""
	if ok {
		flashMessage = fmt.Sprintf("%v", storedFlash)
		delete(session.Values, "flash")

		if err := session.Save(r, w); err != nil {
			app.serverError(w, err)
		}
	}

	td.Flash = flashMessage
	td.IsAuthenticated = app.isAuthenticated(r)
	td.CSRFTemplateTag = csrf.TemplateField(r)

	return td
}

func (app *application) isAuthenticated(r *http.Request) bool {
	session, err := app.session.Get(r, "session-name")
	if err != nil {
		return false
	}

	return session.Values["authenticatedUserID"] != nil
}
