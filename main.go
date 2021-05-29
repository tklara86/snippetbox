package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var infoLog *log.Logger



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


func createSnippet (w http.ResponseWriter, r *http.Request) {
	infoLog = log.New(os.Stdout, "INFO: ", log.LstdFlags)

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

	infoLog.Printf("Create a snippet")

	if err != nil {
		fmt.Println(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {

	infoLog = log.New(os.Stdout, "INFO: " , log.LstdFlags)

	if r.URL.Path != "/" {
		fmt.Println(r.URL.Path)
		http.NotFound(w, r)
		return
	}

	_, err := fmt.Fprintf(w, "<h1>Hello from Home Page</h1>")

	infoLog.Printf("Home page")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	srv := http.Server{
		Addr: ":4000",
		Handler: mux,
	}

	log.Println("Starting server on 4000")
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalln(err)
	}
}
