package main

import (
	"context"
	"flag"
	"github.com/tklara86/snippetbox/cmd/config"
	"github.com/tklara86/snippetbox/cmd/handlers"
	"github.com/tklara86/snippetbox/cmd/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)



func main() {

	app := &config.AppConfig{
		InfoLog: log.New(os.Stdout, "INFO - ", log.LstdFlags),
		ErrorLog: log.New(os.Stderr, "ERROR - ", log.LstdFlags | log.Lshortfile),
	}

	addr := flag.String("addr", ":8080", "HTTP network address")

	flag.Parse()


	// router - servemux
	sm := http.NewServeMux()
	sm.HandleFunc("/", handlers.Home(app))
	sm.HandleFunc("/snippet", handlers.ShowSnippet(app))
	sm.HandleFunc("/snippet/create", handlers.CreateSnippet(app))

	// creates file server which serves files out the './ui/static'
	fileServer := http.FileServer(http.Dir("./ui/static"))

	sm.Handle("/static/", http.StripPrefix("/static", middleware.Neuter(fileServer)))

	// go run ./cmd/web -addr=":4000"
	srv := http.Server{
		Addr: *addr,
		Handler: sm,
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
