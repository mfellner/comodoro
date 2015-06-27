package fleet

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/mfellner/comodoro/api"
	"github.com/mfellner/comodoro/db"
	"github.com/mfellner/comodoro/model"
)

var bucketName = []byte("units")

// CreateUnit creates a new fleet unit.
func CreateUnit(db *db.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var unit model.Unit
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
		if err := json.Unmarshal(body, &unit); err != nil {
			api.BadRequest(w, err.Error())
			return
		}

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			v := b.Get([]byte(unit.Name))

			if v != nil {
				api.DuplicateError(w)
				return err
			}

			jsonString, err := json.Marshal(unit.Body)

			if err != nil {
				api.BadRequest(w, err.Error())
				return err
			}

			if err := b.Put([]byte(unit.Name), []byte(jsonString)); err != nil {
				api.ServerError(w, err.Error())
			} else {
				api.Created(w, unit)
			}
			return err
		})
	})
}

// GetUnits returns a collection of all fleet units.
func GetUnits(db *db.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		units := model.Units{}

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			b.ForEach(func(k, v []byte) error {

				var body map[string]interface{}

				if err := json.Unmarshal(v, &body); err != nil {
					api.ServerError(w, err.Error())
				}

				units = append(units, model.Unit{Name: string(k), Body: body})
				return nil
			})
			return nil
		})

		api.JSON(w, units)
	})
}

// GetUnit returns a single fleet unit for the given name.
func GetUnit(db *db.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		name := vars["name"]

		log.Print(vars)

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			v := b.Get([]byte(name))

			if v == nil {
				api.NotFound(w)
				return nil
			}

			if body, err := unmarshal(v); err != nil {
				api.ServerError(w, err.Error())
			} else {
				api.JSON(w, model.Unit{Name: name, Body: body})
			}

			return nil
		})
	})
}

// DeleteUnit deletes the fleet unit with the given name.
func DeleteUnit(db *db.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		name := vars["name"]

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			v := b.Get([]byte(name))

			if v == nil {
				api.NotFound(w)
				return nil
			}

			if err := b.Delete([]byte(name)); err != nil {
				api.ServerError(w, err.Error())
				return err
			}
			api.Deleted(w)
			return nil
		})
	})
}

func unmarshal(v []byte) (map[string]interface{}, error) {
	var body map[string]interface{}
	err := json.Unmarshal(v, &body)
	return body, err
}
