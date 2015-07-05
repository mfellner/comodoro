package app

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/mfellner/comodoro/rest"
)

var bucketName = []byte("units")

// CreateUnit creates a new fleet unit.
func CreateUnit(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var unit FleetUnit
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
		if err := json.Unmarshal(body, &unit); err != nil {
			rest.BadRequest(w, err.Error())
			return
		}

		app.DB().Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			v := b.Get([]byte(unit.Name))

			if v != nil {
				rest.DuplicateError(w)
				return err
			}

			jsonString, err := json.Marshal(unit.Body)

			if err != nil {
				rest.BadRequest(w, err.Error())
				return err
			}

			if err := b.Put([]byte(unit.Name), []byte(jsonString)); err != nil {
				rest.ServerError(w, err.Error())
			} else {
				rest.Created(w, unit)
			}
			return err
		})
	})
}

// GetUnits returns a collection of all fleet units.
func GetUnits(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		units := FleetUnits{}

		app.DB().View(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			b.ForEach(func(k, v []byte) error {

				var body map[string]interface{}

				if err := json.Unmarshal(v, &body); err != nil {
					rest.ServerError(w, err.Error())
				}

				units = append(units, FleetUnit{Name: string(k), Body: body})
				return nil
			})
			return nil
		})

		rest.JSON(w, units)
	})
}

// GetUnit returns a single fleet unit for the given name.
func GetUnit(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		name := vars["name"]

		app.DB().View(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			v := b.Get([]byte(name))

			if v == nil {
				rest.NotFound(w)
				return nil
			}

			if body, err := unmarshal(v); err != nil {
				rest.ServerError(w, err.Error())
			} else {
				rest.JSON(w, FleetUnit{Name: name, Body: body})
			}

			return nil
		})
	})
}

// DeleteUnit deletes the fleet unit with the given name.
func DeleteUnit(app *App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		name := vars["name"]

		app.DB().Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			v := b.Get([]byte(name))

			if v == nil {
				rest.NotFound(w)
				return nil
			}

			if err := b.Delete([]byte(name)); err != nil {
				rest.ServerError(w, err.Error())
				return err
			}
			rest.Deleted(w)
			return nil
		})
	})
}

func unmarshal(v []byte) (map[string]interface{}, error) {
	var body map[string]interface{}
	err := json.Unmarshal(v, &body)
	return body, err
}
