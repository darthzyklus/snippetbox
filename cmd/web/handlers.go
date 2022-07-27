package main

import (
	"darthzyklus/snippetbox/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

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
	data := app.newTemplateData(req)

	app.render(writer, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(writer http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()

	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	title := req.PostForm.Get("title")
	content := req.PostForm.Get("content")

	expires, err := strconv.Atoi(req.PostForm.Get("expires"))

	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	fieldErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"

	} else if utf8.RuneCountInString(title) > 100 {
		fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "this field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldErrors["expires"] = "This field must be equal to 1, 7 or 365"
	}

	if len(fieldErrors) > 0 {
		fmt.Fprint(writer, fieldErrors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
