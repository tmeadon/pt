package data

import (
	"time"

	"gorm.io/gorm"
)

type ExerciseHistory struct {
	Base
	ExerciseId int           `json:"-"`
	Exercise   Exercise      `json:"exercise"`
	UserId     int           `json:"-"`
	User       User          `json:"user"`
	Sets       []ExerciseSet `json:"sets"`
	WorkoutId  int           `json:"-"`
}

func NewExerciseHistory(userId int, workoutId int, exerciseId int) *ExerciseHistory {
	now := time.Now().UTC()
	return &ExerciseHistory{
		Base: Base{
			Created:      now,
			LastModified: now,
		},
		ExerciseId: exerciseId,
		UserId:     userId,
		WorkoutId:  workoutId,
	}
}

func (eh *ExerciseHistory) AddSet(set *ExerciseSet) {
	eh.LastModified = time.Now().UTC()
	eh.Sets = append(eh.Sets, *set)
}

func (db DB) GetExerciseHistory(id int) (*ExerciseHistory, error) {
	var exerciseHistory ExerciseHistory
	result := db.gorm.Preload("Exercise").Preload("User").Preload("Sets").Find(&exerciseHistory, id)
	return &exerciseHistory, interpretError(result.Error)
}

func (db DB) SaveExerciseHistory(history *ExerciseHistory) error {
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
		if err := tx.Save(history).Error; err != nil {
			return err
		}

		workout.AddExerciseCategory(&exercise.Category)
		if err := tx.Save(workout).Error; err != nil {
			return err
		}

		return nil
	})

	return interpretError(err)
}

func (db DB) DeleteExerciseHistory(history *ExerciseHistory) error {
	history.LastModified = time.Now().UTC()
	history.IsDeleted = true
	result := db.gorm.Save(history)
	db.updateWorkoutCategories(&Workout{Base: Base{Id: history.WorkoutId}})
	return interpretError(result.Error)
}
