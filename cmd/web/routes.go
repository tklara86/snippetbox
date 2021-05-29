package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	// mux
	mux := http.NewServeMux()

	// handlers
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// serve static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))

	return mux

}