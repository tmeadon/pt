package data

import (
	"time"
)

type Workout struct {
	Base
	Name               string             `json:"name"`
	UserId             int                `json:"user_id"`
	User               User               `json:"user"`
	ExerciseInstances  []ExerciseHistory  `json:"exercises"`
	ExerciseCategories []ExerciseCategory `json:"exercise_categories" gorm:"many2many:workout_categories"`
}

func NewWorkout(userID int) *Workout {
	now := time.Now().UTC()
	return &Workout{
		Base: Base{
			Created:      now,
			LastModified: now,
		},
		UserId: userID,
	}
}

func (w *Workout) SetExerciseCategories() {
	for _, e := range w.ExerciseInstances {
		if !workoutContainsCategory(w, &e.Exercise.Category) {
			w.AddExerciseCategory(&e.Exercise.Category)
		}
	}
}

func (w *Workout) AddExerciseCategory(cat *ExerciseCategory) {
	w.LastModified = time.Now().UTC()
	w.ExerciseCategories = append(w.ExerciseCategories, *cat)
}

func workoutContainsCategory(workout *Workout, cat *ExerciseCategory) bool {
	for _, c := range workout.ExerciseCategories {
		if cat.Id == c.Id {
			return true
		}
	}
	return false
}

func (db DB) GetAllWorkouts() ([]Workout, error) {
	var workout []Workout
	result := db.gorm.Where("is_deleted = false").Preload("User").Preload("ExerciseCategories").Find(&workout, "is_deleted == false")
	return workout, interpretError(result.Error)
}

func (db DB) GetWorkout(id int) (*Workout, error) {
	var workout Workout
	result := db.gorm.Where("is_deleted = false").Preload("User").Preload("ExerciseCategories").Preload("ExerciseInstances", "is_deleted = false").Preload("ExerciseInstances.Exercise").Preload("ExerciseInstances.Sets").First(&workout, id)
	return &workout, interpretError(result.Error)
}

func (db DB) SaveWorkout(w *Workout) error {
	w.LastModified = time.Now().UTC()
	result := db.gorm.Save(w)
	db.updateWorkoutCategories(w)
	return interpretError(result.Error)
}

func (db DB) DeleteWorkout(workout *Workout) error {
	workout.LastModified = time.Now().UTC()
	workout.IsDeleted = true
	result := db.gorm.Save(workout)
	return interpretError(result.Error)
}

func (db DB) updateWorkoutCategories(workout *Workout) {
	var w Workout
	db.gorm.Preload("ExerciseInstances", "is_deleted = false").Preload("ExerciseInstances.Exercise.Category").First(&w, workout.Id)

	cats := make([]ExerciseCategory, 0)
	for _, e := range w.ExerciseInstances {
		cats = append(cats, e.Exercise.Category)
	}

	db.gorm.Model(&w).Association("ExerciseCategories").Replace(cats)
}
