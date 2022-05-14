package db

const musclesTableCreate string = `
  CREATE TABLE IF NOT EXISTS muscles (
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  isFront BOOLEAN
);`

func createMusclesTable() error {
	_, err := DB.Exec(musclesTableCreate)
	return err
}
