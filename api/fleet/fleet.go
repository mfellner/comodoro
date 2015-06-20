package fleet

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/mfellner/comodoro/api"
	"github.com/mfellner/comodoro/db"
	"github.com/mfellner/comodoro/model"
)

var BUCKET_NAME_UNITS = []byte("units")

func CreateUnit(db *db.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var unit model.Unit
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
		if err := json.Unmarshal(body, &unit); err != nil {
			api.BadRequest(w, err)
		}

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(BUCKET_NAME_UNITS)
			v := b.Get([]byte(unit.Name))

			if v != nil {
				api.DuplicateError(w)
				return err
			}

			if err := b.Put([]byte(unit.Name), []byte(unit.Body)); err != nil {
				api.ServerError(w, err)
			} else {
				api.Created(w, unit)
			}
			return err
		})
	}
}

func GetUnits(db *db.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		units := model.Units{}

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(BUCKET_NAME_UNITS)
			b.ForEach(func(k, v []byte) error {
				units = append(units, model.Unit{Name: string(k), Body: string(v)})
				return nil
			})
			return nil
		})

		api.WriteResponse(w, units)
	}
}

func GetUnit(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		name := vars["name"]

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(BUCKET_NAME_UNITS)
			v := b.Get([]byte(name))

			if v != nil {
				api.WriteResponse(w, model.Unit{Name: name, Body: string(v)})
			} else {
				api.NotFound(w)
			}
			return nil
		})
	}
}
