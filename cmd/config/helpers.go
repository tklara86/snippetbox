package config

import (
	"bytes"
	"fmt"
	"github.com/tklara86/snippetbox/pkg/models"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"
)
type TemplateData struct {
	CurrentYear 	int
	FormData  		url.Values   // url.Values, which is the same underlying type as the r.PostForm map that held the data sent in the request body.
	FormErrors  	map[string]string
	Snippet 		*models.Snippet
	Snippets		[]*models.Snippet
}

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

func (app *AppConfig) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	if td == nil {
		td = &TemplateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}

func (app *AppConfig) Render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {

	// Retrieve the appropriate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.
	ts, ok := app.TemplateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("template %s does not exist", name))
	}

	// Initialize buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError helper and then
	// return.

	err := ts.Execute(buf, app.addDefaultData(td,r) )
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter. Again, this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Fatal(err)
	}

}


