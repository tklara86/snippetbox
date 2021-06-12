package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/tklara86/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)


type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *mysql.SnippetModel
	users		  *mysql.UserModel
	session 	  *sessions.Session
	templateCache map[string]*template.Template
}

func main() {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	sessionSecret := os.Getenv("SECRET")


	// address flag
	addr := flag.String("addr", ":4000", "HTTP network address")

	// define DSN string
	dsn := flag.String("dsn", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName), "MySQL data source name")

	// secret
	secret := flag.String("secret", sessionSecret, "secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// init new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		session: session,
		snippets: &mysql.SnippetModel{
			DB: db,
		},
		templateCache: templateCache,
		users: &mysql.UserModel{
			DB: db,
		},
	}

	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}
	colorizeTerminalMsg(ColorBlue)
	infoLog.Printf("Started server on port %s", srv.Addr)
	err = srv.ListenAndServe()

	if err != nil {
		colorizeTerminalMsg(ColorRed)
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	// opens database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Verifies connection to the database is still alive
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
