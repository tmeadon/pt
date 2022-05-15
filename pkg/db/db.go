package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var dbPath string

func ConnectDatabase(path string) error {
	dbPath = path
	ensureDbExists(dbPath)

	db, err := sql.Open("sqlite3", (dbPath + "?_foreign_keys=true"))
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

func recordExists(id int, tableName string) (bool, error) {
	stmt, err := DB.Prepare(fmt.Sprintf("SELECT id from %s WHERE id = ?", tableName))
	if err != nil {
		return false, err
	}

	err = stmt.QueryRow(id).Scan(&id)

	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}

	return true, nil
}
