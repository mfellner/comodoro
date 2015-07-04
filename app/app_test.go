package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func withApp(fn func(*App)) func() {
	return func() {
		var db DB
		if err := db.Open("/tmp/comodoro.db", 0600); err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		defer os.Remove("/tmp/comodoro.db")

		fn(NewApp(&db))
	}
}

func withFleetMock(fn func()) func() {
	return func() {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, "foobar")
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		defer server.Close()

		viper.SetDefault("fleetEndpoint", server.URL)
		fn()
	}
}

// TestApp tests the application.
func TestApp(t *testing.T) {

	convey.Convey("Given a fleet unit JSON", t, withFleetMock(withApp(func(app *App) {

		router := app.Router()

		unit := map[string]interface{}{
			"name": "test-unit",
			"body": map[string]interface{}{
				"foo": "bar",
			},
		}

		jsonString, err := json.Marshal(unit)
		if err != nil {
			log.Fatal(err)
		}

		convey.Convey("When the unit is first submitted", func() {
			req := newPOSTRequest("/api/fleet/units", jsonString)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			convey.Convey("Then the response should be status Created", func() {
				convey.So(w.Code, convey.ShouldEqual, 201)
			})

			convey.Convey("And when it is fetched again", func() {
				req := newGETRequest(fmt.Sprintf("/api/fleet/units/%s", "test-unit"))
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				convey.Convey("Then the response should be status OK", func() {
					convey.So(w.Code, convey.ShouldEqual, 200)
				})
			})
		})

		convey.Convey("When the same unit is submitted twice", func() {
			req := newPOSTRequest("/api/fleet/units", jsonString)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			convey.Convey("The API should return status Created on the first request", func() {
				convey.So(w.Code, convey.ShouldEqual, 201)
			})

			req = newPOSTRequest("/api/fleet/units", jsonString)
			w = httptest.NewRecorder()

			router.ServeHTTP(w, req)

			convey.Convey("The API should return error Duplicate Entry on the second request", func() {
				convey.So(w.Code, convey.ShouldEqual, 409)
			})
		})

	})))

	convey.Convey("Given the Comodoro App", t, withApp(func(app *App) {

		convey.Convey("When a non-existing unit is requested", func() {
			req := newGETRequest("/api/fleet/units/foobar")
			w := httptest.NewRecorder()
			app.Router().ServeHTTP(w, req)

			convey.Convey("Then it should return error Not Found", func() {
				convey.So(w.Code, convey.ShouldEqual, 404)
			})
		})

		convey.Convey("When I send an OPTIONS request", func() {
			req := newOPTIONSRequest("/api/fleet/units")
			req.Header.Add("Origin", "http://localhost")
			w := httptest.NewRecorder()
			app.Router().ServeHTTP(w, req)

			convey.Convey("Then it should return status OK", func() {
				convey.So(w.Code, convey.ShouldEqual, 200)
				convey.So(w.Header().Get("Access-Control-Allow-Origin"), convey.ShouldEqual, "http://localhost")
				convey.So(w.Header().Get("Access-Control-Allow-Headers"), convey.ShouldEqual, "Content-Type")
				convey.So(w.Header().Get("Access-Control-Allow-Methods"), convey.ShouldEqual, "POST, GET, PUT, DELETE, OPTIONS")
				convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "text/plain")
			})
		})
	}))
}

func newPOSTRequest(url string, json []byte) *http.Request {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func newGETRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func newOPTIONSRequest(url string) *http.Request {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	return req
}
