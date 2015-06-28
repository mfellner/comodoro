package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// TestApi tests the API.
func TestApi(t *testing.T) {

	convey.Convey("Given the Index handler", t, func() {
		handler := Index()

		convey.Convey("When I send a GET request", func() {
			req := newGETRequest()
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			convey.Convey("Then then it should return status 200 OK", func() {
				convey.So(w.Code, convey.ShouldEqual, 200)
				convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})
	})

	convey.Convey("Given an object", t, func() {
		unit := map[string]interface{}{
			"name": "test-unit",
			"body": map[string]interface{}{
				"foo": "bar",
			},
		}

		convey.Convey("When I pass it to the JSON handler", func() {
			w := httptest.NewRecorder()
			JSON(w, unit)

			convey.Convey("Then then it should return status 200 OK", func() {
				convey.So(w.Code, convey.ShouldEqual, 200)
				convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})

		convey.Convey("When I pass an unsupported type to the JSON handler", func() {
			callJSON := func() {
				w := httptest.NewRecorder()
				JSON(w, func() {})
				log.Print(w.Body)
			}
			convey.Convey("Then then it should panic", func() {
				convey.So(callJSON, convey.ShouldPanic)
			})
		})

		convey.Convey("When I pass it to the Created handler", func() {
			w := httptest.NewRecorder()
			Created(w, unit)

			convey.Convey("Then then it should return status 201 OK", func() {
				convey.So(w.Code, convey.ShouldEqual, 201)
				convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "application/json; charset=UTF-8")
			})
		})
	})

	convey.Convey("Given an error message", t, func() {
		message := "Test Error"

		convey.Convey("When I pass it to the BadRequest handler", func() {
			w := httptest.NewRecorder()
			BadRequest(w, message)

			convey.Convey("Then then it should return status 400 Bad Request", func() {
				convey.So(w.Code, convey.ShouldEqual, 400)
				convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "text/plain")
			})
		})

		convey.Convey("When I pass it to the ServerError handler", func() {
			w := httptest.NewRecorder()
			ServerError(w, message)

			convey.Convey("Then then it should return status 500 Internal Error", func() {
				convey.So(w.Code, convey.ShouldEqual, 500)
				convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "text/plain")
			})
		})
	})

	convey.Convey("When I call the Deleted handler", t, func() {
		w := httptest.NewRecorder()
		Deleted(w)

		convey.Convey("Then then it should return status 204 No Content", func() {
			convey.So(w.Code, convey.ShouldEqual, 204)
			convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "")
		})
	})

	convey.Convey("When I call the Not Found handler", t, func() {
		w := httptest.NewRecorder()
		NotFound(w)

		convey.Convey("Then then it should return status 404 Not Found", func() {
			convey.So(w.Code, convey.ShouldEqual, 404)
			convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "text/plain")
		})
	})

	convey.Convey("When I call the Duplicate Error handler", t, func() {
		w := httptest.NewRecorder()
		DuplicateError(w)

		convey.Convey("Then then it should return status 409 Conflict", func() {
			convey.So(w.Code, convey.ShouldEqual, 409)
			convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "text/plain")
		})
	})
}

func newGETRequest() *http.Request {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Fatal(err)
	}
	return req
}
