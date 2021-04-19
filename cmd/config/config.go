package config

import (
	"github.com/golangcollege/sessions"
	"github.com/tklara86/snippetbox/pkg/models/postgres"
	"html/template"
	"log"
)

type AppConfig struct {
	InfoLog  		*log.Logger
	ErrorLog 		*log.Logger
	Session 		*sessions.Session
	Snippets 		*postgres.SnippetModel
	TemplateCache    map[string]*template.Template
}
