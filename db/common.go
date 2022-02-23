package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	DEFAULT_DB_PORT = 5432

	STATUS_ACTIVE   = "A"
	STATUS_DELETED  = "D"
	STATUS_DISABLED = "S"

	ROLE_ADMIN      = "A"
	ROLE_POWERADMIN = "P"
	ROLE_USER       = "U"
	ROLE_SERVICE    = "S"

	POWERDOMAIN = 1 // Domain_id for powerdomain (superuser)
)

var DB_HOST = "localhost"
var DB_PORT = 5432
var DB_USER = "postgres"
var DB_PASSWORD = "glowglow"
var DB_NAME = "unidb"
var DB_DISABLE_SSL = false

// Function for handling errors
func checkErr(err error) {
	if err == nil {
		return
	}

	log.Printf("err=%v\n", err)
	panic(err)
}

func init() {
	var str string

	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Error loading .env file")
	}

	str = os.Getenv("DB_HOST")
	if str != "" {
		DB_HOST = str
	}

	str = os.Getenv("DB_PORT")
	if str != "" {
		DB_PORT, _ = strconv.Atoi(str)
	}

	str = os.Getenv("DB_USER")
	if str != "" {
		DB_USER = str
	}

	str = os.Getenv("DB_PASSWORD")
	if str != "" {
		DB_PASSWORD = str
	}

	str = os.Getenv("DB_NAME")
	if str != "" {
		DB_NAME = str
	}

	str = os.Getenv("DB_DISABLE_SSL")
	if str != "" {
		DB_DISABLE_SSL, _ = strconv.ParseBool(str)
	}
	setupDB()
}

var dbHandle *sql.DB

// DB set up
func setupDB() *sql.DB {
	if dbHandle == nil {
		dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
			DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
		if DB_DISABLE_SSL {
			dbinfo += " sslmode=disable"
		}

		dbHandle, _ = sql.Open("postgres", dbinfo)
		ctx, stop := context.WithCancel(context.Background())
		defer stop()
		if err := dbHandle.PingContext(ctx); err != nil {
			log.Fatalf("unable to connect to database: %v", err)
		} else {
			log.Printf("DB: %s Successful\n", dbinfo)
		}
	}

	return dbHandle
}

// func resetDB() {
// 	closeDB()
// }

// func closeDB() {
// 	if dbHandle != nil {
// 		dbHandle.Close()
// 		dbHandle = nil
// 	}
// }
