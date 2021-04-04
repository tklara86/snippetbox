package main

import (
	"context"
	"flag"
	"github.com/tklara86/snippetbox/cmd/handlers"
	"github.com/tklara86/snippetbox/cmd/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)


func main() {
	// Log errors
	errorLog := log.New(os.Stderr, "ERROR - ", log.LstdFlags | log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO - ", log.LstdFlags)

	addr := flag.String("addr", ":8080", "HTTP network address")

	flag.Parse()

	// router - servemux
	sm := http.NewServeMux()
	sm.HandleFunc("/", handlers.Home)
	sm.HandleFunc("/snippet", handlers.ShowSnippet)
	sm.HandleFunc("/snippet/create", handlers.CreateSnippet)

	// creates file server which serves files out the './ui/static'
	fileServer := http.FileServer(http.Dir("./ui/static"))

	sm.Handle("/static/", http.StripPrefix("/static", middleware.Neuter(fileServer)))

	// go run ./cmd/web -addr=":4000"
	srv := http.Server{
		Addr: *addr,
		Handler: sm,
		ErrorLog: errorLog,
	}
	infoLog.Printf("Starting server on %s", *addr)
	// start the server
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			errorLog.Fatal(err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	errorLog.Println("Received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
}
