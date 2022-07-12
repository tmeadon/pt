package data

import (
	"time"

	"gorm.io/gorm"
)

func (db DB) GetAllExercises() ([]Exercise, error) {
	var exercises []Exercise
	result := db.gorm.Preload("Category").Preload("Muscles").Preload("SecondaryMuscles").Preload("Equipment").Find(&exercises)
	return exercises, interpretError(result.Error)
}

func (db DB) GetExercise(id int) (*Exercise, error) {
	var exercise Exercise
	result := db.gorm.Preload("Category").Preload("Muscles").Preload("SecondaryMuscles").Preload("Equipment").First(&exercise, id)
	return &exercise, interpretError(result.Error)
}

func (db DB) GetExerciseHistory(id int) (*ExerciseHistory, error) {
	var exerciseHistory ExerciseHistory
	result := db.gorm.Preload("Exercise").Preload("User").Preload("Sets").Find(&exerciseHistory, id)
	return &exerciseHistory, interpretError(result.Error)
}

func (db DB) InsertExercise(exercise *Exercise) error {
	exercise.LastModified = time.Now().UTC()
	exercise.Created = time.Now().UTC()
	result := db.gorm.Create(exercise)
	return interpretError(result.Error)
}

func (db DB) InsertExerciseHistory(history *ExerciseHistory) error {
	history.Created = time.Now().UTC()
	history.LastModified = time.Now().UTC()

	workout, err := db.GetWorkout(history.WorkoutId)
	if err != nil {
		return interpretError(err)
	}

	exercise, err := db.GetExercise(history.ExerciseId)
	if err != nil {
		return interpretError(err)
	}

	err = db.gorm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(history).Error; err != nil {
			return err
		}

		if !workoutContainsCategory(workout, &exercise.Category) {
			workout.ExerciseCategories = append(workout.ExerciseCategories, exercise.Category)
			workout.LastModified = time.Now().UTC()
			if err := tx.Save(&workout).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return interpretError(err)
}

func (db DB) UpdateExerciseHistory(history *ExerciseHistory) error {
	history.LastModified = time.Now().UTC()
	result := db.gorm.Save(history)
	return interpretError(result.Error)
}
