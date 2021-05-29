package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)


// home Home page
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		fmt.Println(r.URL.Path)
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/footer.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w,nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}

}


// showSnippet displays snippet with specific id
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	_, err = fmt.Fprintf(w,"Snippet with ID %d", id)

	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "Fprintf: %v\n", err)

		if err != nil {
			fmt.Println(err)
		}
	}

}
// createSnippet creates a new snippet
func (app *application) createSnippet (w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_,err := w.Write([]byte("Method not allowed"))

		if err != nil {
			app.clientError(w, http.StatusMethodNotAllowed)
		}
		return
	}
	_, err := fmt.Fprintf(w,"<h1>Create snippet route</h1>")

	if err != nil {
		fmt.Println(err)
	}
}


