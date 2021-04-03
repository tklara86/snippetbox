package handlers

import (
	"fmt"
	"net/http"
	"strconv"
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
	// Get the id parameter from url e.g /snippet?id=123
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	_, err = fmt.Fprintf(w, "Display a specific with ID %d", id)
	if err != nil {
		http.Error(w, "could not get a snippet", http.StatusBadRequest)
	}
}

// CreateSnippet handler function
func CreateSnippet(w http.ResponseWriter, r *http.Request) {

	// Check if HTTP method POST
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