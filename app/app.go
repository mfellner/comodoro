package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfellner/comodoro/db"
)

type App struct {
	db     *db.DB
	router *mux.Router
}

func NewApp(d *db.DB) *App {

	app := &App{
		db:     d,
		router: mux.NewRouter().StrictSlash(true),
	}

	for _, route := range routes(app.db) {
		app.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return app
}

func (a *App) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, a.router)
}
