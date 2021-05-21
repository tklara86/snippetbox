package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello form Home page"))

	if err != nil {
		fmt.Println(err)
	}
}


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	err := http.ListenAndServe(":4000", mux)

	if err != nil {
		log.Fatalln(err)
	}
}
