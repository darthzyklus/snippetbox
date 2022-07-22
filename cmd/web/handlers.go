package main

import (
	"darthzyklus/snippetbox/internal/models"
	"errors"
	"fmt"

	"net/http"
	"strconv"
)

func (app *application) home(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		app.notFound(writer)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(writer, err)
		return
	}

	/*
		for _, snippet := range snippets {
			fmt.Fprintf(writer, "%+v\n", snippet)
		}
	*/
	page := "home.tmpl"

	data := app.newTemplateData(req)
	data.Snippets = snippets

	app.render(writer, 200, page, data)
}

func (app *application) snippetView(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}

	page := "view.tmpl"

	data := app.newTemplateData(req)
	data.Snippet = snippet

	app.render(writer, 200, page, data)
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
