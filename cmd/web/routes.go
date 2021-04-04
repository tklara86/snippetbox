package main

import (
	"github.com/tklara86/snippetbox/cmd/config"
	"github.com/tklara86/snippetbox/cmd/handlers"
	"github.com/tklara86/snippetbox/cmd/middleware"
	"net/http"
)

func routes(app *config.AppConfig) *http.ServeMux {
	sm := http.NewServeMux()
	sm.Handle("/", handlers.Home(app))
	sm.Handle("/snippet", handlers.ShowSnippet(app))
	sm.Handle("/snippet/create", handlers.CreateSnippet(app))

	// creates file server which serves files out the './ui/static'
	fileServer := http.FileServer(http.Dir("./ui/static"))

	sm.Handle("/static/", http.StripPrefix("/static", middleware.Neuter(fileServer)))

	return sm
}
