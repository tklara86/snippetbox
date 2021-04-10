package config

import (
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"
)

func (app *AppConfig) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *AppConfig) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *AppConfig) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

func (app *AppConfig) Render(w http.ResponseWriter, r *http.Request, name string, td *template.Template) {
	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.
	ts, ok := app.TemplateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("The template %s does not exist", name))
	}
	// Execute the template set, passing in any dynamic data.
	err := ts.Execute(w, td)
	if err != nil {
		app.ServerError(w, err)
	}

}
