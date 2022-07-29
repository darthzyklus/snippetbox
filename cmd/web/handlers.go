package main

import (
	"darthzyklus/snippetbox/internal/models"
	"darthzyklus/snippetbox/internal/validator"
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
	data := app.newTemplateData(req)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(writer, http.StatusOK, "create.tmpl", data)
}

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
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
		Title:   req.PostForm.Get("title"),
		Content: req.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more thant 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must be equal to 1,7 or 365")

	if !form.Valid() {
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
