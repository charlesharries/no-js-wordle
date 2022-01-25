package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charlesharries/wordle-go/pkg/forms"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	g, _ := app.session.Get(r, "game").(game)
	if len(g.Answer) == 0 {
		g.Answer = app.randomWord()
		app.session.Put(r, "game", g)
	}

	attempt := r.URL.Query().Get("attempt")

	app.render(w, r, "home.page.tmpl", &templateData{
		Attempt: attempt,
		Form:    forms.New(nil),
		Game:    g,
	})
}

func (app *application) guess(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	g, ok := app.session.Get(r, "game").(game)
	if !ok {
		g = game{}
	}

	form := forms.New(r.PostForm)
	fmt.Println(form.Get("attempt"))
	form.Required("attempt")
	form.MaxLength("attempt", 5)
	form.MinLength("attempt", 5)
	form.InWordList("attempt", app.words)

	if !form.Valid() {
		app.render(w, r, "home.page.tmpl", &templateData{
			Form: form,
			Game: g,
		})
		return
	}

	g.Guess(strings.ToLower(form.Get("attempt")))
	app.session.Put(r, "game", g)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) restart(w http.ResponseWriter, r *http.Request) {
	blank := game{}
	app.session.Put(r, "game", blank)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) letter(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	g, ok := app.session.Get(r, "game").(game)
	if !ok {
		g = game{}
	}

	form := forms.New(r.PostForm)
	form.Required("letter")

	if !form.Valid() {
		app.render(w, r, "home.page.tmpm", &templateData{
			Form: form,
			Game: g,
		})
	}

	var attempt string
	if form.Get("letter") == "backspace" {
		attempt = trimLastChar(form.Get("attempt"))
	} else {
		attempt = fmt.Sprintf("%s%s", form.Get("attempt"), form.Get("letter"))
	}

	http.Redirect(w, r, fmt.Sprintf("/?attempt=%s", firstN(attempt, 5)), http.StatusSeeOther)
}
