package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.home)))
	mux.Post("/", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.guess)))
	mux.Post("/letter", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.letter)))
	mux.Post("/restart", dynamicMiddleware.ThenFunc(http.HandlerFunc(app.restart)))

	// Serve static files from the /ui/static directory.
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
