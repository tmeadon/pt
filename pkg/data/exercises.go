package data

import "gorm.io/gorm"

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
	result := db.gorm.Create(exercise)
	return interpretError(result.Error)
}

func (db DB) InsertExerciseHistory(history *ExerciseHistory) error {
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
			if err := tx.Save(&workout).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return interpretError(err)
}
