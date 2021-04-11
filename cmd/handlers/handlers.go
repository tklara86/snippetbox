package handlers

import (
	"errors"
	"fmt"
	"github.com/tklara86/snippetbox/cmd/config"
	"github.com/tklara86/snippetbox/pkg/models"
	"net/http"
	"strconv"
	"time"
)

// Home Handler
func Home(app *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			// replies with 404 HTTP status error
			http.NotFound(w, r)
			return
		}

		s, err := app.Snippets.Latest()
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// Render page and pass the data
		app.Render(w,r, "home.page.tmpl", &config.TemplateData{Snippets: s})
	})
}

// ShowSnippet handler
func ShowSnippet(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id parameter from url e.g /snippet?id=123
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			app.NotFound(w) // Use the NotFound() helper
			return
		}
		// Use the SnippetModel object's Get method to retrieve the data for a
		// specific record based on its ID. If no matching record is found,
		// return a 404 Not Found response.
		s, err := app.Snippets.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				app.NotFound(w)
			} else {
				app.ServerError(w, err)
			}
			return
		}

		// Render page and pass the data
		app.Render(w,r, "show.page.tmpl", &config.TemplateData{Snippet: s})
	}

}


// CreateSnippet handler function
func CreateSnippet(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.ClientError(w, http.StatusMethodNotAllowed)
			return
		}

		title := "O snail2"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := time.Now().Add(time.Hour * 24 * 10)

		id, err := app.Snippets.Insert(title,content,expires)

		if err != nil {
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
	}

}