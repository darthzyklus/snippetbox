package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(writer, req)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(writer, "Interal Server Errror", http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(writer, "base", nil)

	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(writer, "Interal Server Errror", http.StatusInternalServerError)
	}
}

func (app *application) snippetView(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(writer, req)
		return
	}

	fmt.Fprintf(writer, "Display a specific snippet with ID %d", id)
}

func (app *application) snippetCreate(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		http.Error(writer, "Method Not Allowed", http.StatusUnsupportedMediaType)
		return
	}

	writer.Write([]byte("Create a snippet..."))
}
