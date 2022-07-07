package main

import (
	"log"
	"net/http"
)

func home(writer http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		http.NotFound(writer, req)
		return
	}

	writer.Write([]byte("Hello from Snippetbox"))
}

func snippetView(writer http.ResponseWriter, req *http.Request) {
	writer.Write([]byte("Display a specific snippet..."))
}

func snippetCreate(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		writer.Header().Set("Allow", "POST")
		writer.WriteHeader(405)
		writer.Write([]byte("Method Not Allowed"))
		return
	}

	writer.Write([]byte("Create a new snippet..."))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("Starting server on port 4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
