package main

import (
	"log"
	"net/http"
)

func home(writer http.ResponseWriter, req *http.Request) {
	writer.Write([]byte("Hello from Snippetbox"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Println("Starting server on port 4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
