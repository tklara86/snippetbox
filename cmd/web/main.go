package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
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

	flag.Parse()
	// go run ./cmd/web -addr=":4000"

	db, err := openDB(*dsn)


	// Initialize a mysql.SnippetModel instance and add it to the application
	// dependencies.
	app := &config.AppConfig{
		InfoLog: log.New(os.Stdout, "INFO - ", log.LstdFlags),
		ErrorLog: log.New(os.Stderr, "ERROR - ", log.LstdFlags | log.Lshortfile),
		Snippets: &postgres.SnippetModel{
			DB: db,
		},
	}

	if err != nil {
		app.ErrorLog.Fatal(err)
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
	}
	app.InfoLog.Printf("Starting server on port %s", *addr)

	// start the server
	go func() {
		err := srv.ListenAndServe()
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
