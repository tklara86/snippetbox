package handlers

import (
	"fmt"
	"github.com/tklara86/snippetbox/cmd/config"
	"html/template"
	"net/http"
	"strconv"
)

// Home Handler

func Home (app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// replies with 404 HTTP status error
			http.NotFound(w, r)
			return
		}

		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/footer.partial.tmpl",
		}

		// Parses tmpl files
		tmpl, err := template.ParseFiles(files...)
		if err != nil {

			app.ErrorLog.Println(err.Error())
			http.Error(w, "internal Server error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Printf(err.Error())
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
func ShowSnippet(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id parameter from url e.g /snippet?id=123
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			http.NotFound(w, r)
			return
		}
		_, err = fmt.Fprintf(w, "Display a specific with ID %d", id)
		if err != nil {
			http.Error(w, "could not get a snippet", http.StatusBadRequest)
		}
	}
}


// CreateSnippet handler function
func CreateSnippet(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if HTTP method POST
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(405)
			_, err := w.Write([]byte("Method Not Allowed"))

			if err != nil {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		_, err := w.Write([]byte(`{"msg": "Create a new snippet"}`))
		if err != nil {
			http.Error(w, "could not create a snippet", http.StatusBadRequest)
		}
	}

}