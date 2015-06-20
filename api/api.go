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

func WriteResponse(w http.ResponseWriter, v interface{}) {
	writeResponse(w, http.StatusOK, v)
}

func Index(w http.ResponseWriter, r *http.Request) {
	v := struct{}{}
	WriteResponse(w, v)
}

func Created(w http.ResponseWriter, v interface{}) {
	writeResponse(w, http.StatusCreated, v)
}

func BadRequest(w http.ResponseWriter, v interface{}) {
	writeResponse(w, http.StatusBadRequest, v)
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "Not Found")
}

func DuplicateError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusConflict)
	fmt.Fprint(w, "Duplicate Entry")
}

func ServerError(w http.ResponseWriter, v interface{}) {
	writeResponse(w, http.StatusInternalServerError, v)
}
