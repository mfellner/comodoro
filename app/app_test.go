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

		unitName := "test-unit"

		unit := map[string]interface{}{
			"name": unitName,
			"body": map[string]interface{}{
				"foo": "bar",
			},
		}

		jsonString, err := json.Marshal(unit)
		if err != nil {
			log.Fatal(err)
		}

		convey.Convey("When the unit is first submitted", func() {
			w, req := newPOSTRequest("/api/fleet/units", jsonString)
			router.ServeHTTP(w, req)

			convey.Convey("Then the response should be status Created", func() {
				convey.So(w.Code, convey.ShouldEqual, 201)
			})

			convey.Convey("And when it is fetched again", func() {
				w, req := newGETRequest(fmt.Sprintf("/api/fleet/units/%s", unitName))
				router.ServeHTTP(w, req)

				convey.Convey("Then the response should be that unit", func() {
					actualObj := unmarshalBuffer(w.Body)
					convey.So(w.Code, convey.ShouldEqual, 200)
					convey.So(actualObj, convey.ShouldResemble, unit)
				})
			})

			convey.Convey("And when the list of all units is fetched", func() {
				w, req := newGETRequest("/api/fleet/units")
				router.ServeHTTP(w, req)

				convey.Convey("Then the response should be a list that contains the unit", func() {
					actualObj := unmarshalBuffer(w.Body)
					expectedObj := []interface{}{unit}
					convey.So(w.Code, convey.ShouldEqual, 200)
					convey.So(actualObj, convey.ShouldResemble, expectedObj)
				})
			})

			convey.Convey("And when it is deleted", func() {
				w, req := newDELETERequest(fmt.Sprintf("/api/fleet/units/%s", unitName))
				router.ServeHTTP(w, req)

				convey.Convey("Then the response should be status 204 No Content", func() {
					convey.So(w.Code, convey.ShouldEqual, 204)
				})

				convey.Convey("And when it is fetched again", func() {
					w, req := newGETRequest(fmt.Sprintf("/api/fleet/units/%s", unitName))
					router.ServeHTTP(w, req)

					convey.Convey("Then the response should be status 404 Not Found", func() {
						convey.So(w.Code, convey.ShouldEqual, 404)
					})
				})
			})
		})

		convey.Convey("When the same unit is submitted twice", func() {
			w, req := newPOSTRequest("/api/fleet/units", jsonString)
			router.ServeHTTP(w, req)

			convey.Convey("The API should return status Created on the first request", func() {
				convey.So(w.Code, convey.ShouldEqual, 201)
			})

			w, req = newPOSTRequest("/api/fleet/units", jsonString)
			router.ServeHTTP(w, req)

			convey.Convey("The API should return error Duplicate Entry on the second request", func() {
				convey.So(w.Code, convey.ShouldEqual, 409)
			})
		})

	})))

	convey.Convey("Given the Comodoro App", t, withApp(func(app *App) {

		convey.Convey("When a non-existing unit is requested", func() {
			w, req := newGETRequest("/api/fleet/units/foobar")
			app.Router().ServeHTTP(w, req)

			convey.Convey("Then it should return error Not Found", func() {
				convey.So(w.Code, convey.ShouldEqual, 404)
			})
		})

		convey.Convey("When I send an OPTIONS request", func() {
			w, req := newOPTIONSRequest("/api/fleet/units")
			req.Header.Add("Origin", "http://localhost")
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

func unmarshalBuffer(b *bytes.Buffer) interface{} {
	var obj interface{}
	if err := json.Unmarshal(b.Bytes(), &obj); err != nil {
		log.Fatal(err)
	}
	return obj
}

func newRequest(method string, url string, json []byte) (*httptest.ResponseRecorder, *http.Request) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(json))
	if err != nil {
		log.Fatal(err)
	}
	return httptest.NewRecorder(), req
}

func newPOSTRequest(url string, json []byte) (*httptest.ResponseRecorder, *http.Request) {
	return newRequest("POST", url, json)
}

func newGETRequest(url string) (*httptest.ResponseRecorder, *http.Request) {
	return newRequest("GET", url, nil)
}

func newDELETERequest(url string) (*httptest.ResponseRecorder, *http.Request) {
	return newRequest("DELETE", url, nil)
}

func newOPTIONSRequest(url string) (*httptest.ResponseRecorder, *http.Request) {
	return newRequest("OPTIONS", url, nil)
}
