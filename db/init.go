package db

import (
	"log"
	"os"
	"sync"

	"go.etcd.io/bbolt"
)

var (
	DB    *bbolt.DB
	Mutex sync.Mutex
)

// Initialize the database
func InitDB() {
	var err error
	DB, err = bbolt.Open("mydb.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("persons"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Close the database
func CloseDB() {
	DB.Close()
	// os.Remove("mydb.db") // Uncomment this line to remove the database file after execution
}
