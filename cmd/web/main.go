package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog *log.Logger
	errorLog *log.Logger
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

	// address flag
	addr := flag.String("addr", ":4000", "HTTP network address")

	// define DSN string
	dsn := flag.String("dsn", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName), "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags |log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		infoLog: infoLog,
		errorLog: errorLog,
	}

	srv := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		ErrorLog: errorLog,
	}
	colorizeTerminalMsg(ColorGreen)
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
