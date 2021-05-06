package handlers

import (
	"errors"
	"fmt"
	"github.com/tklara86/snippetbox/cmd/config"
	"github.com/tklara86/snippetbox/pkg/forms"
	"github.com/tklara86/snippetbox/pkg/models"
	"log"
	"net/http"
	"strconv"
	"time"
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
		app.Render(w,r, "home.page.tmpl", &config.TemplateData{
			Snippets: s,
		})
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
		app.Render(w,r,"create.page.tmpl", &config.TemplateData{
			// Pass a new empty forms.Form object to the template.
			Form: forms.New(nil),
		})
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
		// Create a new forms.Form struct containing the POSTed data from the
		// form, then use the validation methods to check the content.
		form := forms.New(r.PostForm)
		form.Required("title", "content", "expires")
		form.MaxLength("title", 100)
	//	form.PermittedValues("expires", "365", "7", "1")

		// If the form isn't valid, redisplay the template passing in the
		// form.Form object as the data.
		if !form.Valid() {
			app.Render(w, r, "create.page.tmpl", &config.TemplateData{Form: form})
			return
		}

		// Because the form data (with type url.Values) has been anonymously embedded
		// in the form.Form struct, we can use the Get() method to retrieve
		// the validated value for a particular form field.
		expiresString := form.Get("expires")

		i, _ := strconv.Atoi(expiresString)
		exp := time.Duration(i)
		expires := time.Now().Add(time.Hour * 24 * exp)

		id, err := app.Snippets.Insert(form.Get("title"),form.Get("content"), expires)

		if err != nil {
			app.ServerError(w, err)
			return
		}

		app.Session.Put(r, "flash", "Snippet successfully created")

		http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
	}

}

// SignupUserForm signs up a new user
func SignupUserForm(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		app.Render(w,r, "signup.page.tmpl", &config.TemplateData{
			Form: forms.New(nil),
		})
	}
}

// SignUpUser creates a new user
func SignUpUser(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			app.ClientError(w, http.StatusBadRequest)
			return
		}
		// Validate the form contents
		form := forms.New(r.PostForm)
		form.Required("name", "email", "password")
		form.MaxLength("name", 255)
		form.MaxLength("email", 255)
		form.MatchesPattern("email", forms.EmailRx)
		form.MinLength("password", 10)

		if !form.Valid() {
			app.Render(w,r,"signup.page.tmpl", &config.TemplateData{
				Form: form,
			})
			return
		}
		// Otherwise send a placeholder response (for now!).
		err = app.Users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))

		if err != nil {
			app.ServerError(w, err)
			return
		}
		app.Session.Put(r, "flash", "Your signup was successful. Please log in.")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)

	}



}

// LoginUserForm displays the user login form
func LoginUserForm(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Display the user login form...")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// LoginUser authenticates and login the user
func LoginUser(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Authenticate and login the user..")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// LogoutUser logouts the user
func LogoutUser(app *config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Logout the user...")
		if err != nil {
			log.Fatalln(err)
		}
	}
}