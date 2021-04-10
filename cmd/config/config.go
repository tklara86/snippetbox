package config

import (
	"github.com/tklara86/snippetbox/pkg/models/postgres"
	"log"
)

type AppConfig struct {
	InfoLog  		*log.Logger
	ErrorLog 		*log.Logger
	Snippets 		*postgres.SnippetModel
}
