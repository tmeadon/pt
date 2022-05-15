package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
)

func Backup() (string, error) {
	var conns []*sqlite3.SQLiteConn

	// register a custom driver that provides access to the underlying sqlite3.SQLiteConn type
	sql.Register("sqlite3_backup",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(sc *sqlite3.SQLiteConn) error {
				conns = append(conns, sc)
				return nil
			}},
	)

	// connect to the src
	srcDb, err := sql.Open("sqlite3_backup", dbPath)
	if err != nil {
		return "", err
	}
	defer srcDb.Close()
	srcDb.Ping()

	// create and connect to backup db
	dstPath := dbPath + "-" + time.Now().Format("2006-01-02T15:04:05-0700")
	err = ensureDbExists(dstPath)
	if err != nil {
		return "", err
	}

	dstDb, err := sql.Open("sqlite3_backup", dstPath)
	if err != nil {
		return "", err
	}
	defer dstDb.Close()
	dstDb.Ping()

	bk, err := conns[1].Backup("main", conns[0], "main")
	if err != nil {
		return "", err
	}

	_, err = bk.Step(-1)
	if err != nil {
		return "", err
	}

	err = bk.Finish()
	fmt.Printf("Backed up %s to %s\n", dbPath, dstPath)
	return dstPath, err
}
