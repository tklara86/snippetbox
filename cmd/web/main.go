package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog *log.Logger
	errorLog *log.Logger
}

func main() {
	// address flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags |log.Lshortfile)

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
	err := srv.ListenAndServe()

	if err != nil {
		colorizeTerminalMsg(ColorRed)
		errorLog.Fatal(err)
	}
}
