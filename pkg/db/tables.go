package db

import "fmt"

type tableMap struct {
	musclesTable string
}

var tables = tableMap{
	musclesTable: "muscles",
}

const dropTableString string = "DROP TABLE IF EXISTS %s;"

func RecreateMusclesTable() error {
	err := dropTable(tables.musclesTable)
	if err != nil {
		return err
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS %s (
      id INTEGER NOT NULL PRIMARY KEY,
      name TEXT NOT NULL,
      simple_name TEXT NOT NULL,
      is_front BOOLEAN
    );`
	_, err = DB.Exec(fmt.Sprintf(createTable, tables.musclesTable))

	return err
}

func dropTable(table string) error {
	_, err := DB.Exec(fmt.Sprintf(dropTableString, table))
	return err
}
