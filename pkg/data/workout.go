package data

import (
	"time"
)

func (db DB) GetAllWorkouts() ([]Workout, error) {
	var workout []Workout
	result := db.gorm.Preload("User").Find(&workout, "is_deleted == false")
	return workout, interpretError(result.Error)
}

func (db DB) GetWorkout(id int) (*Workout, error) {
	var workout Workout
	result := db.gorm.Preload("User").Preload("ExerciseCategories").Preload("ExerciseInstances.Exercise").Preload("ExerciseInstances.Sets").First(&workout, id)
	return &workout, interpretError(result.Error)
}

func (db DB) InsertWorkout(workout *Workout) error {
	workout.Created = time.Now().UTC()
	workout.LastModified = time.Now().UTC()
	result := db.gorm.Create(workout)
	return interpretError(result.Error)
}

func (db DB) UpdateWorkout(workout *Workout) error {
	workout.LastModified = time.Now().UTC()
	result := db.gorm.Save(workout)
	return interpretError(result.Error)
}

func workoutContainsCategory(workout *Workout, cat *ExerciseCategory) bool {
	for _, c := range workout.ExerciseCategories {
		if cat.Id == c.Id {
			return true
		}
	}
	return false
}

func (db DB) DeleteWorkout(workout *Workout) error {
	workout.LastModified = time.Now().UTC()
	workout.IsDeleted = true
	result := db.gorm.Save(workout)
	return interpretError(result.Error)
}
