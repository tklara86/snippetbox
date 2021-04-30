package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"github.com/tklara86/snippetbox/cmd/config"
	"github.com/tklara86/snippetbox/cmd/handlers"

	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.RecoverPanic, app.LogRequest, config.SecureHeaders)

	dynamicMiddleware := alice.New(app.Session.Enable)


	sm := pat.New()
	sm.Get("/", dynamicMiddleware.Then(handlers.Home(app)))
	sm.Get("/snippet/create", dynamicMiddleware.Then(handlers.CreateSnippetForm(app)))
	sm.Post("/snippet/create", dynamicMiddleware.Then(handlers.CreateSnippet(app)))
	sm.Get("/snippet/:id", dynamicMiddleware.Then(handlers.ShowSnippet(app)))

	sm.Get("/user/signup", dynamicMiddleware.Then(handlers.SignupUserForm(app)))
	sm.Post("/user/signup", dynamicMiddleware.Then(handlers.SignUpUser(app)))
	sm.Get("/user/login", dynamicMiddleware.Then(handlers.LoginUserForm(app)))
	sm.Post("/user/login", dynamicMiddleware.Then(handlers.LoginUser(app)))
	sm.Post("/user/logout", dynamicMiddleware.Then(handlers.LogoutUser(app)))

	// creates file server which serves files out the './ui/static'
	fileServer := http.FileServer(http.Dir("./ui/static"))

	sm.Get("/static/", http.StripPrefix("/static", config.Neuter(fileServer)))

	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(sm)
}
