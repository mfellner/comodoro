package app

import (
	"net/http"

	"github.com/mfellner/comodoro/api"
	"github.com/mfellner/comodoro/api/fleet"
	"github.com/mfellner/comodoro/db"
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

func routes(d *db.DB) Routes {
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
			fleet.CreateUnit(d),
		},
		Route{
			"GetFleetUnits",
			"GET",
			"/api/fleet/units",
			fleet.GetUnits(d),
		},
		Route{
			"GetFleetUnit",
			"GET",
			"/api/fleet/units/{name}",
			fleet.GetUnit(d),
		},
		Route{
			"GetFleetUnit",
			"DELETE",
			"/api/fleet/units/{name}",
			fleet.DeleteUnit(d),
		},
	}
}
