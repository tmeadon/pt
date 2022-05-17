package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbPath string

func ConnectDatabase(path string) error {
	dbPath = path

	d, err := sql.Open("sqlite3", (dbPath + "?_foreign_keys=true"))
	if err != nil {
		return err
	}

	err = d.Ping()
	if err != nil {
		return err
	}

	db = d

	return nil
}

func queryRow(query string, queryArgs ...any) (*sql.Row, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(queryArgs...)
	return row, nil
}

func queryRows(query string, queryArgs ...any) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(queryArgs...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return rows, nil
}
