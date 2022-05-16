package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var dbPath string

func ConnectDatabase(path string) error {
	dbPath = path

	db, err := sql.Open("sqlite3", (dbPath + "?_foreign_keys=true"))
	if err != nil {
		return err
	}

	DB = db

	return nil
}
