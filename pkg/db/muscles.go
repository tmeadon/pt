package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tmeadon/pt/pkg/models"
)

func InsertMuscle(newMuscle models.Muscle, keepID bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(getMuscleInsertStatement(keepID))
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = execMuscleInsert(stmt, &newMuscle, keepID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func getMuscleInsertStatement(keepID bool) string {
	if keepID {
		return fmt.Sprintf("INSERT INTO %s (id, name, simple_name, is_front) VALUES (?, ?, ?, ?)", tables.MusclesTable)
	}
	return fmt.Sprintf("INSERT INTO %s (name, simple_name, is_front) VALUES (?, ?, ?)", tables.MusclesTable)
}

func execMuscleInsert(stmt *sql.Stmt, newMuscle *models.Muscle, keepID bool) error {
	if keepID {
		_, err := stmt.Exec(newMuscle.Id, newMuscle.Name, newMuscle.SimpleName, newMuscle.IsFront)
		return err
	}
	_, err := stmt.Exec(newMuscle.Name, newMuscle.SimpleName, newMuscle.IsFront)
	return err
}

func GetAllMuscles() ([]models.Muscle, error) {
	muscles := make([]models.Muscle, 0)

	query := "select id, name, simple_name, is_front from %s"
	rows, err := queryRows(fmt.Sprintf(query, tables.MusclesTable))

	if rows == nil {
		return muscles, err
	}

	for rows.Next() {
		muscle := new(models.Muscle)
		err := rows.Scan(&muscle.Id, &muscle.Name, &muscle.SimpleName, &muscle.IsFront)
		if err != nil {
			return []models.Muscle{}, err
		}
		muscles = append(muscles, *muscle)
	}

	return muscles, nil
}

func GetMuscle(id int) DBResponse {
	query := `select m.id, m.name, m.simple_name, m.is_front
	from %s as m
	where m.id = ?
	`

	row, err := queryRow(fmt.Sprintf(query, tables.MusclesTable), id)
	if err != nil {
		return DBResponse{Status: Fail, Error: err}
	}

	muscle := models.Muscle{}
	err = row.Scan(&muscle.Id, &muscle.Name, &muscle.SimpleName, &muscle.IsFront)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DBResponse{Status: Empty}
		}
		return DBResponse{Status: Fail, Error: err}
	}

	return DBResponse{Status: Success, Data: muscle}
}
