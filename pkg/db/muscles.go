package db

import (
	"database/sql"
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

	stmt, err := db.Prepare(fmt.Sprintf("SELECT id, name, simple_name, is_front from %s", tables.MusclesTable))
	if err != nil {
		return []models.Muscle{}, err
	}

	rows, err := stmt.Query()
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Muscle{}, nil
		}
		return []models.Muscle{}, err
	}
	defer rows.Close()

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
