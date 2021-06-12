package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler{

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, securityHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)
	// mux
	mux := pat.New()

	// handlers
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// serve static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", neuter(fileServer)))

	return standardMiddleware.Then(mux)

}

