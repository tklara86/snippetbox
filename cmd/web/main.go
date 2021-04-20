package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	"github.com/tklara86/snippetbox/cmd/config"
	"github.com/tklara86/snippetbox/pkg/models/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
)


// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}


func main() {


	// Load env variables
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading env file")
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	dbname := os.Getenv("DBNAME")




	addr := flag.String("addr", ":8080", "HTTP network address")

	dsn := flag.String("dsn", "postgres://" + username + ":" + password + "@" + host + "/" + dbname + "?sslmode=disable", "Postgres data source name")

	// Define a new command-line flag for the session secret (a random key which
	// will be used to encrypt and authenticate session cookies). It should be 32
	// bytes long.
	secret := flag.String("secret", "9a7d0e35b825c535550d63b1bb95a8e099764feb53a3eb9bc3976fb430ea70af", "Secret key")

	flag.Parse()
	// go run ./cmd/web -addr=":4000"

	db, err := openDB(*dsn)

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		log.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	//Initialize a tls.Config struct to hold the non-default TLS settings we want
	// the server to use.
	tlsConfig := tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}


	// And add the session manager to our application dependencies.
	app := &config.AppConfig{
		InfoLog: log.New(os.Stdout, "INFO - ", log.LstdFlags),
		ErrorLog: log.New(os.Stderr, "ERROR - ", log.LstdFlags | log.Lshortfile),
		Session: session,
		Snippets: &postgres.SnippetModel{
			DB: db,
		},
		TemplateCache: templateCache,
	}


	// Defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits
	defer func() {
		if err := db.Close(); err != nil {
			app.ErrorLog.Fatal(err)
		}
	}()



	srv := http.Server{
		Addr: *addr,
		Handler: routes(app),
		ErrorLog: app.ErrorLog,
		TLSConfig: &tlsConfig,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	app.InfoLog.Printf("Starting server on port %s", *addr)

	// start the server
	go func() {
		err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		app.ErrorLog.Fatal(err)

	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	app.InfoLog.Println("Received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
}
