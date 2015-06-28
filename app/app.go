package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfellner/comodoro/db"
	"github.com/spf13/viper"
)

// App encapsulates the resources of the web service.
type App struct {
	db          *db.DB
	router      *mux.Router
	fleetClient *FleetClient
}

// NewApp creates and configures a new application instance.
func NewApp(d *db.DB) *App {
	fleetEndpoint := viper.GetString("fleetEndpoint")

	app := &App{
		db:          d,
		router:      mux.NewRouter().StrictSlash(true),
		fleetClient: NewFleetClient(fleetEndpoint),
	}

	for _, route := range routes(app) {
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
