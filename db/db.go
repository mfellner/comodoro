package db

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

// DB encapsulates the database resources.
type DB struct {
	*bolt.DB
}

// Open and initialize the database.
func (db *DB) Open(path string, mode os.FileMode) error {
	var err error

	db.DB, err = bolt.Open(path, mode, nil)
	if err != nil {
		return err
	}

	// Create 'unit' bucket.
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("units"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		db.Close()
		return err
	}

	return nil
}
