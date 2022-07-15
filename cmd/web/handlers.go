package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		app.notFound(writer)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(writer, err)
		return
	}

	err = ts.ExecuteTemplate(writer, "base", nil)

	if err != nil {
		app.serverError(writer, err)
	}
}

func (app *application) snippetView(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	fmt.Fprintf(writer, "Display a specific snippet with ID %d", id)
}

func (app *application) snippetCreate(writer http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	title := "GO AND RUST"
	content := "Best snippets ever\n Find here the best content\n\n Andres Uris"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, req, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
