package main

import (
	"flag"
	"fmt"

	"log"
	"net/http"
	"os"
)

const (
	ColorBlack   	  = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"
)

func colorizeTerminalMsg(color string) {
	fmt.Println(color)
}

func main() {
	// address flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR\t", log.LstdFlags |log.Lshortfile)

	// mux
	mux := http.NewServeMux()
	// handlers
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// serve static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))

	srv := &http.Server{
		Addr: *addr,
		Handler: mux,
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
