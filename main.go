package main

import (
	"log"
	"net/http"
)


// Home Handler
func Home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello form SnippetBox"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func main() {
	sm := http.NewServeMux()
	sm.HandleFunc("/", Home)

	srv := http.Server{
		Addr: ":8080",
		Handler: sm,
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
