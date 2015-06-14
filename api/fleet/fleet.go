package fleet

import (
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/mfellner/comodoro/api"
	"github.com/mfellner/comodoro/db"
	"github.com/mfellner/comodoro/model"
)

func GetUnits(db *db.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		units := model.Units{}

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("units"))
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
			b := tx.Bucket([]byte("units"))
			v := b.Get([]byte(name))

			if v != nil {
				unit := model.Unit{
					Name: name,
					Body: string(v),
				}
				api.WriteResponse(w, unit)
			} else {
				api.NotFound(w, r)
			}

			return nil
		})
	}
}
