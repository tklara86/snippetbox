package main

import (
	"context"
	"flag"
	"github.com/tklara86/snippetbox/cmd/config"
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
	// go run ./cmd/web -addr=":4000"

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
