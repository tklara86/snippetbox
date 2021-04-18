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

	sm := pat.New()
	sm.Get("/", handlers.Home(app))
	sm.Get("/snippet/create", handlers.CreateSnippetForm(app))
	sm.Post("/snippet/create", handlers.CreateSnippet(app))
	sm.Get("/snippet/:id", handlers.ShowSnippet(app))

	// creates file server which serves files out the './ui/static'
	fileServer := http.FileServer(http.Dir("./ui/static"))

	sm.Get("/static/", http.StripPrefix("/static", config.Neuter(fileServer)))

	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(sm)
}
