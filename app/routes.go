package app

import (
	"net/http"

	"github.com/mfellner/comodoro/api"
)

// Route represents a route for the HTTP handler.
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

// Routes is a collection of routes.
type Routes []Route

func routes(app *App) Routes {
	return Routes{
		Route{
			"Index",
			"GET",
			"/",
			api.Index(),
		},
		Route{
			"CreateFleetUnit",
			"POST",
			"/api/fleet/units",
			CreateUnit(app),
		},
		Route{
			"GetFleetUnits",
			"GET",
			"/api/fleet/units",
			GetUnits(app),
		},
		Route{
			"GetFleetUnit",
			"GET",
			"/api/fleet/units/{name}",
			GetUnit(app),
		},
		Route{
			"GetFleetUnit",
			"DELETE",
			"/api/fleet/units/{name}",
			DeleteUnit(app),
		},
	}
}
