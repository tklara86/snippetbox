package main

import (
	"context"
	"github.com/tklara86/snippetbox/cmd/handlers"
	"github.com/tklara86/snippetbox/cmd/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)


func main() {
	l := log.New(os.Stdout, "snippetbox-", log.LstdFlags)
	// router - servemux
	sm := http.NewServeMux()
	sm.HandleFunc("/", handlers.Home)
	sm.HandleFunc("/snippet", handlers.ShowSnippet)
	sm.HandleFunc("/snippet/create", handlers.CreateSnippet)

	// creates file server which serves files out the './ui/static'
	fileServer := http.FileServer(http.Dir("./ui/static"))

	sm.Handle("/static/", http.StripPrefix("/static", middleware.Neuter(fileServer)))

	srv := http.Server{
		Addr: ":8080",
		Handler: sm,
		ErrorLog: l,
	}

	// start the server
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = srv.Shutdown(ctx)
}
