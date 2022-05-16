package db

import (
	"bytes"
	"fmt"
	"html/template"
)

var schemaVersionTemplates = map[int][]string{
	1: v1_command_templates,
}

func ApplyMigrations() error {
	currentSchema, _ := getCurrentSchemaVersion()

	for k, v := range schemaVersionTemplates {
		if k > currentSchema {
			err := executeMigration(&v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getCurrentSchemaVersion() (int, error) {
	stmt, err := DB.Prepare(fmt.Sprintf("SELECT version FROM %s WHERE id = 0", tables.SchemaVersionTable))
	if err != nil {
		return 0, err
	}

	var version int = 0
	err = stmt.QueryRow().Scan(&version)
	if err != nil {
		return 0, err
	}

	return version, nil
}

func executeMigration(commandTemplates *[]string) error {
	for _, t := range *commandTemplates {
		cmd, err := buildTableCommand(t)
		if err != nil {
			return err
		}

		stmt, err := DB.Prepare(cmd)
		if err != nil {
			return err
		}

		_, err = stmt.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func buildTableCommand(tableTemplate string) (cmd string, err error) {
	tmpl, err := template.New("schema").Parse(tableTemplate)
	if err != nil {
		return
	}

	var schemaBytes bytes.Buffer
	err = tmpl.Execute(&schemaBytes, tables)
	if err != nil {
		return
	}

	cmd = schemaBytes.String()
	return
}

var v1_command_templates []string = []string{
	`CREATE TABLE IF NOT EXISTS {{.SchemaVersionTable}} (
		id INTEGER NOT NULL PRIMARY KEY,
		version INTEGER NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS {{.MusclesTable}} (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		simple_name TEXT NOT NULL,
		is_front BOOLEAN
	);`,
	`CREATE TABLE IF NOT EXISTS {{.EquipmentTable}} (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS {{.ExerciseCategoriesTable}} (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL
	);`,
	`CREATE TABLE IF NOT EXISTS {{.ExercisesTable}} (
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		category_id INTEGER NOT NULL,
		FOREIGN KEY (category_id) REFERENCES {{.ExerciseCategoriesTable}} (id)
	);`,
	`CREATE TABLE IF NOT EXISTS {{.ExerciseEquipmentTable}} (
		id INTEGER NOT NULL PRIMARY KEY,
		exercise_id INTEGER NOT NULL,
		equipment_id INTEGER NOT NULL,
		FOREIGN KEY (exercise_id) REFERENCES {{.ExercisesTable}} (id),
		FOREIGN KEY (equipment_id) REFERENCES {{.EquipmentTable}} (id)
	);`,
	`CREATE TABLE IF NOT EXISTS {{.ExerciseMusclesTable}} (
		id INTEGER NOT NULL PRIMARY KEY,
		exercise_id INTEGER NOT NULL,
		muscle_id INTEGER NOT NULL,
		is_primary BOOLEAN,
		FOREIGN KEY (exercise_id) REFERENCES {{.ExercisesTable}} (id),
		FOREIGN KEY (muscle_id) REFERENCES {{.MusclesTable}} (id)
	);`,
	`INSERT INTO {{.SchemaVersionTable}} (id, version) VALUES (0, 1);`,
}
