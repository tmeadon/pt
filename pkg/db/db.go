package db

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var dbPath string = "./pt.db"
var DB *sql.DB

func ConnectDatabase() error {
	ensureDbExists(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	DB = db

	return nil
}

func ensureDbExists(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(path)
		return err
	}
	return nil
}

func Init() error {
	return createMusclesTable()
}
