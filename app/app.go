package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfellner/comodoro/db"
)

// App encapsulates the resources of the web service.
type App struct {
	db     *db.DB
	router *mux.Router
}

// NewApp creates and configures a new application instance.
func NewApp(d *db.DB) *App {

	app := &App{
		db:     d,
		router: mux.NewRouter().StrictSlash(true),
	}

	for _, route := range routes(app.db) {
		app.router.
			Methods(route.Method, "OPTIONS").
			Path(route.Pattern).
			Name(route.Name).
			Handler(HandleCORS(AllowOrigin(route.Handler)))
	}

	return app
}

// DB is the application's database.
func (a *App) DB() *db.DB {
	return a.db
}

// Router is the application's router.
func (a *App) Router() *mux.Router {
	return a.router
}

// ListenAndServe starts the application on the given port.
func (a *App) ListenAndServe(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), a.router)
}
