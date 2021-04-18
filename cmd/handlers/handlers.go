package handlers

import (
	"errors"
	"fmt"
	"github.com/tklara86/snippetbox/cmd/config"
	"github.com/tklara86/snippetbox/pkg/models"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Home Handler
func Home(app *config.AppConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
		// Pat doesn't strip the colon from the named capture key, so we need to
		// get the value of ":id" from the query string instead of "id".
		id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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

func CreateSnippetForm(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		app.Render(w,r,"create.page.tmpl", nil)
	}
}


// CreateSnippet handler function
func CreateSnippet(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// First we call r.ParseForm() which adds any data in POST request bodies
		// to the r.PostForm map. This also works in the same way for PUT and PATCH
		// requests. If there are any errors, we use our app.ClientError helper to send
		// a 400 Bad Request response to the user.
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}
		// Use the r.PostForm.Get() method to retrieve the relevant data fields
		// from the r.PostForm map.
		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")
		expiresString := r.PostForm.Get("expires")


		// Initialize a mpa to hold any validation errors
		errors := make(map[string]string)

		// Check that the Title field isn't blank.
		if strings.TrimSpace(title) == "" {
			errors["title"] = "This field cannot be blank"
		}
		// Check that the Title field does not have more than 100 characters
		if utf8.RuneCountInString(title) > 100 {
			errors["title"] = "This field is too long (maximum is 100 characters)"
		}

		// Check that the Content field isn't blank.
		if strings.TrimSpace(content) == "" {
			errors["content"] = "This field cannot be blank"
		}

		// Check that the Expires field isn't blank.
		if strings.TrimSpace(expiresString) == "" {
			errors["expires"] = "This field cannot be blank"
		}

		//Check the expires field matches one of the permitted
		// values ("1", "7" or "365").
		if expiresString != "365" && expiresString != "7" && expiresString != "1" {
			errors["expires"] = "This field is invalid"
		}

		if len(errors) > 0 {
			app.Render(w, r, "create.page.tmpl", &config.TemplateData{
				FormData: r.PostForm,
				FormErrors: errors,
			})
			return
		}
		i, _ := strconv.Atoi(expiresString)
		exp := time.Duration(i)
		expires := time.Now().Add(time.Hour * 24 * exp)

		id, err := app.Snippets.Insert(title, content, expires)

		if err != nil {
			app.ServerError(w, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	}

}