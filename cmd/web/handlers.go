package main

import (
	"darthzyklus/snippetbox/internal/models"
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(writer http.ResponseWriter, req *http.Request) {

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(writer, err)
		return
	}

	page := "home.tmpl"

	data := app.newTemplateData(req)
	data.Snippets = snippets

	app.render(writer, 200, page, data)
}

func (app *application) snippetView(writer http.ResponseWriter, req *http.Request) {
	params := httprouter.ParamsFromContext(req.Context())

	id, err := strconv.Atoi(params.ByName("id"))

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
	writer.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(writer http.ResponseWriter, req *http.Request) {

	title := "GO AND RUST"
	content := "Best snippets ever\n Find here the best content\n\n Andres Uris"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
