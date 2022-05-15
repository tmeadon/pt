package db

import (
	"database/sql"
	"fmt"

	"github.com/tmeadon/pt/pkg/models"
)

func InsertExercise(newExercise models.Exercise, keepID bool) error {
	// start a transaction
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// insert the exercise
	err = insertExercise(tx, &newExercise, keepID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// insert primary and secondary exercise muscle links
	err = insertExerciseMuscles(tx, newExercise.Id, &newExercise.Muscles, true)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = insertExerciseMuscles(tx, newExercise.Id, &newExercise.MusclesSecondary, false)
	if err != nil {
		tx.Rollback()
		return err
	}

	// insert equipment links
	err = insertAllExerciseEquipment(tx, newExercise.Id, &newExercise.Equipment)
	if err != nil {
		tx.Rollback()
		return err
	}

	// finally commit the transaction
	err = tx.Commit()
	return err
}

func insertExercise(tx *sql.Tx, newEx *models.Exercise, keepId bool) error {
	stmt, err := tx.Prepare(getExerciseInsertStatement(keepId))
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = execExerciseInsert(stmt, newEx, keepId)
	return err
}

func getExerciseInsertStatement(keepID bool) string {
	if keepID {
		return fmt.Sprintf("INSERT INTO %s (id, name, description, category_id) VALUES (?, ?, ?, ?)", tables.ExercisesTable)
	}
	return fmt.Sprintf("INSERT INTO %s (name, description, category_id) VALUES (?, ?, ?)", tables.ExercisesTable)
}

func execExerciseInsert(stmt *sql.Stmt, newEx *models.Exercise, keepID bool) error {
	if keepID {
		_, err := stmt.Exec(newEx.Id, newEx.Name, newEx.Description, newEx.Category.Id)
		return err
	}
	_, err := stmt.Exec(newEx.Name, newEx.Description, newEx.Category.Id)
	return err
}

func insertExerciseMuscles(tx *sql.Tx, exerciseId int, muscles *[]models.Muscle, isPrimary bool) error {
	for _, m := range *muscles {
		err := insertExerciseMuscle(tx, exerciseId, m.Id, isPrimary)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertExerciseMuscle(tx *sql.Tx, exerciseId int, muscleId int, isPrimary bool) error {
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (exercise_id, muscle_id, is_primary) VALUES (?, ?, ?)", tables.ExerciseMusclesTable))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(exerciseId, muscleId, isPrimary)
	return err
}

func insertAllExerciseEquipment(tx *sql.Tx, exerciseId int, equipment *[]models.Equipment) error {
	for _, e := range *equipment {
		err := insertExerciseEquipment(tx, exerciseId, e.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertExerciseEquipment(tx *sql.Tx, exerciseId int, equipmentId int) error {
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (exercise_id, equipment_id) VALUES (?, ?)", tables.ExerciseEquipmentTable))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(exerciseId, equipmentId)
	return err
}
