package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func writeResponse(w http.ResponseWriter, i int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(i)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

// Index writes a response for the root of the API.
func Index(w http.ResponseWriter, r *http.Request) {
	v := struct{}{}
	JSON(w, v)
}

// JSON writes a code 200 JSON content response.
func JSON(w http.ResponseWriter, v interface{}) {
	writeResponse(w, http.StatusOK, v)
}

// Created writes a code 201 response for a created entity.
func Created(w http.ResponseWriter, v interface{}) {
	writeResponse(w, http.StatusCreated, v)
}

// BadRequest writes a code 400 response for a client error.
func BadRequest(w http.ResponseWriter, v interface{}) {
	writeResponse(w, http.StatusBadRequest, v)
}

// NotFound writes a code 404 response for a non-existing resource.
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "Not Found")
}

//DuplicateError writes a code 409 response for conflict error.
func DuplicateError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusConflict)
	fmt.Fprint(w, "Duplicate Entry")
}

// ServerError writes a code 500 response for a server-side error.
func ServerError(w http.ResponseWriter, v interface{}) {
	writeResponse(w, http.StatusInternalServerError, v)
}
