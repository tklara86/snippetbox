package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)


// Home Handler
func Home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		// replies with 404 HTTP status error
		http.NotFound(w, r)
		return
	}

	_, err := w.Write([]byte("Hello from SnippetBox"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

// ShowSnippet handler function
func ShowSnippet(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Display specific snippets"))
	if err != nil {
		http.Error(w, "could not get a snippet", http.StatusBadRequest)
	}
}
// CreateSnippet handler function
func CreateSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(405)
		_, err := w.Write([]byte("Method Not Allowed"))

		if err != nil {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	_, err := w.Write([]byte(`{"msg": "Create a new snippet"}`))
	if err != nil {
		http.Error(w, "could not create a snippet", http.StatusBadRequest)
	}
}

func main() {
	l := log.New(os.Stdout, "snippetbox-", log.LstdFlags)
	// router - servemux
	sm := http.NewServeMux()
	sm.HandleFunc("/", Home)
	sm.HandleFunc("/snippet", ShowSnippet)
	sm.HandleFunc("/snippet/create", CreateSnippet)

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
