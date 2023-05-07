package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	once       sync.Once
	instance   *sql.DB
	postgresDB *sql.DB
)

func Get(dbname string) *sql.DB {
	var err error

	switch viper.GetString("env") {
	default:
		once.Do(func() {
			dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
				viper.GetString("database.host"),
				viper.GetInt("database.port"),
				viper.GetString("database.username"),
				viper.GetString("database.password"),
				dbname)
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

	case "test":
		createDatabase(dbname)
		var conn *sql.DB
		dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			viper.GetString("database.host"),
			viper.GetInt("database.port"),
			viper.GetString("database.username"),
			viper.GetString("database.password"),
			dbname)
		conn, err = sql.Open("postgres", dbinfo)
		if err != nil {
			panic(err)
		}

		// test the connection to the database.
		err = conn.Ping()
		if err != nil {
			panic(err)
		}
		return conn
	}

	if err != nil {
		panic(err)
	}

	return instance
}

func createDatabase(dbname string) error {
	once.Do(func() {
		var err error
		dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			viper.GetString("database.host"),
			viper.GetInt("database.port"),
			viper.GetString("database.username"),
			viper.GetString("database.password"),
			"postgres")
		postgresDB, err = sql.Open("postgres", dbinfo)
		if err != nil {
			log.Fatal(err)
		}
	})
	_, err := postgresDB.Exec("CREATE DATABASE " + dbname)
	return err
}
