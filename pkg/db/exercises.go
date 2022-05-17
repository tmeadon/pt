package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tmeadon/pt/pkg/models"
)

func InsertExercise(newExercise models.Exercise, keepID bool) error {
	// start a transaction
	tx, err := db.Begin()
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

func ListAllExercises() ([]models.Exercise, error) {
	// query the db for the exercise bases
	exercises, err := listAllExerciseSummaries()
	return exercises, err
}

func listAllExerciseSummaries() ([]models.Exercise, error) {
	exercises := make([]models.Exercise, 0)

	query := `select e.id, e.name, c.id, c.name 
	from %s as e
	join %s as c on e.category_id = c.id
	`

	rows, err := queryRows(fmt.Sprintf(query, tables.ExercisesTable, tables.ExerciseCategoriesTable))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		e := new(models.Exercise)
		err := rows.Scan(&e.Id, &e.Name, &e.Category.Id, &e.Category.Name)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, *e)
	}

	return exercises, nil
}

func GetExerciseById(id int) (models.Exercise, []error) {
	errs := make([]error, 0)
	exercise, err := getExerciseBase(id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			errs = append(errs, err)
		}
		return models.Exercise{}, errs
	}

	err = getExerciseMuscles(&exercise)
	if err != nil {
		errs = append(errs, err)
	}

	err = getExerciseEquipment(&exercise)
	if err != nil {
		errs = append(errs, err)
	}

	return exercise, errs
}

func getExerciseBase(id int) (models.Exercise, error) {
	query := `select e.id, e.name, e.description, c.id, c.name
	from %s as e
	join %s as c on e.category_id = c.id
	where e.id = ?
	`

	row, err := queryRow(fmt.Sprintf(query, tables.ExercisesTable, tables.ExerciseCategoriesTable), id)
	if err != nil {
		return models.Exercise{}, err
	}

	exercise := models.Exercise{}
	err = row.Scan(&exercise.Id, &exercise.Name, &exercise.Description, &exercise.Category.Id, &exercise.Category.Name)
	if err != nil {
		return models.Exercise{}, err
	}

	return exercise, nil
}

func getExerciseMuscles(exercise *models.Exercise) error {
	query := `select m.id, m.name, m.simple_name, m.is_front, em.is_primary
	from %s as em
	join %s as m on em.muscle_id = m.id
	where em.exercise_id = ?
	`

	rows, err := queryRows(fmt.Sprintf(query, tables.ExerciseMusclesTable, tables.MusclesTable), exercise.Id)
	if err != nil {
		return err
	}

	for rows.Next() {
		var isPrimary bool
		muscle := models.Muscle{}
		err = rows.Scan(&muscle.Id, &muscle.Name, &muscle.SimpleName, &muscle.IsFront, &isPrimary)
		if err != nil {
			return err
		}
		if isPrimary {
			exercise.Muscles = append(exercise.Muscles, muscle)
		} else {
			exercise.MusclesSecondary = append(exercise.MusclesSecondary, muscle)
		}
	}

	return nil
}

func getExerciseEquipment(exercise *models.Exercise) error {
	query := `select e.id, e.name
	from %s as ee
	join %s as e on e.id = ee.equipment_id
	where ee.exercise_id = ?`

	rows, err := queryRows(fmt.Sprintf(query, tables.ExerciseEquipmentTable, tables.EquipmentTable), exercise.Id)
	if err != nil {
		return err
	}

	for rows.Next() {
		equipment := models.Equipment{}
		err = rows.Scan(&equipment.Id, &equipment.Name)
		if err != nil {
			return err
		}
		exercise.Equipment = append(exercise.Equipment, equipment)
	}

	return nil
}
