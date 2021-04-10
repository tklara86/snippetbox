package config

import (
	"github.com/tklara86/snippetbox/pkg/models/postgres"
	"html/template"
	"log"
)

type AppConfig struct {
	InfoLog  		*log.Logger
	ErrorLog 		*log.Logger
	Snippets 		*postgres.SnippetModel
	TemplateCache   map[string]*template.Template
}
