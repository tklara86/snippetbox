package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)


// home Home page
func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		fmt.Println(r.URL.Path)
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/footer.partial.tmpl",
		"./ui/html/base.layout.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w,nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}

}


// showSnippet displays snippet with specific id
func showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	_, err = fmt.Fprintf(w,"Snippet with ID %d", id)

	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "Fprintf: %v\n", err)

		if err != nil {
			fmt.Println(err)
		}
	}

}
// createSnippet creates a new snippet
func createSnippet (w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		_,err := w.Write([]byte("Method not allowed"))

		if err != nil {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	_, err := fmt.Fprintf(w,"<h1>Create snippet route</h1>")



	if err != nil {
		fmt.Println(err)
	}
}


