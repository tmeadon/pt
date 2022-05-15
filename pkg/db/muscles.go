package db

import (
	"database/sql"
	"fmt"

	"github.com/tmeadon/pt/pkg/models"
)

func InsertMuscle(newMuscle models.Muscle) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (name, simple_name, is_front) VALUES (?, ?, ?)", tables.musclesTable))
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newMuscle.Name, newMuscle.SimpleName, newMuscle.IsFront)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func GetAllMuscles() ([]models.Muscle, error) {
	muscles := make([]models.Muscle, 0)

	stmt, err := DB.Prepare(fmt.Sprintf("SELECT id, name, full_name, is_front from %s", tables.musclesTable))
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
