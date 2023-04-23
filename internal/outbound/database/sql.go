package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

var (
	once     sync.Once
	instance *sql.DB
)

func Get() *sql.DB {
	var err error

	// Create the singleton instance of DB
	once.Do(func() {
		dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			"localhost", 5432, "myuser", "mypassword", "mydb")
		instance, err = sql.Open("postgres", dbinfo)
		if err != nil {
			return
		}

		// test the connection to the database.
		err = instance.Ping()
		if err != nil {
			return
		}
	})

	if err != nil {
		panic(err)
	}

	return instance
}
