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

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(writer, http.StatusOK, "create.tmpl", data)
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) snippetCreatePost(writer http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()

	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(req.PostForm.Get("expires"))

	if err != nil {
		app.clientError(writer, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:       req.PostForm.Get("title"),
		Content:     req.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: make(map[string]string),
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"

	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "this field cannot be blank"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must be equal to 1, 7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(req)
		data.Form = form

		app.render(writer, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(writer, err)
		return
	}

	http.Redirect(writer, req, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
