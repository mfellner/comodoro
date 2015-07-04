package middleware

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// TestMiddleware tests the middleware.
func TestMiddleware(t *testing.T) {

	convey.Convey("Given a new http.Handler chain.", t, func() {
		handler := New(HandleCORS, AllowOrigin, LogHTTP).
			Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		convey.Convey("When I send a GET request", func() {
			req := newRequest("GET")
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			convey.Convey("Then then it should return status 200 OK", func() {
				convey.So(w.Code, convey.ShouldEqual, 200)
				convey.So(w.Header().Get("Access-Control-Allow-Origin"), convey.ShouldEqual, "null")
			})
		})

		convey.Convey("When I send an OPTIONS request", func() {
			req := newRequest("OPTIONS")
			req.Header.Add("Origin", "http://localhost")
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			convey.Convey("Then then it should return status 200 OK", func() {
				convey.So(w.Code, convey.ShouldEqual, 200)
				convey.So(w.Header().Get("Access-Control-Allow-Origin"), convey.ShouldEqual, "http://localhost")
				convey.So(w.Header().Get("Access-Control-Allow-Headers"), convey.ShouldEqual, "Content-Type")
				convey.So(w.Header().Get("Access-Control-Allow-Methods"), convey.ShouldEqual, "POST, GET, PUT, DELETE, OPTIONS")
				convey.So(w.Header().Get("Access-Control-Max-Age"), convey.ShouldEqual, "86400")
				convey.So(w.Header().Get("Content-Type"), convey.ShouldEqual, "text/plain")
			})
		})
	})
}

func newRequest(verb string) *http.Request {
	req, err := http.NewRequest(verb, "", nil)
	if err != nil {
		log.Fatal(err)
	}
	return req
}
