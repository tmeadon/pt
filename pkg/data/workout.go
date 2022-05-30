package data

import "time"

func (db DB) GetAllWorkouts() ([]Workout, error) {
	var workout []Workout
	result := db.gorm.Preload("User").Find(&workout)
	return workout, interpretError(result.Error)
}

func (db DB) GetWorkout(id int) (Workout, error) {
	var workout Workout
	result := db.gorm.Preload("User").First(&workout, id)
	return workout, interpretError(result.Error)
}

func (db DB) InsertWorkout(workout *Workout) (*Workout, error) {
	workout.CreatedAt = time.Now().UTC()
	result := db.gorm.Create(workout)
	return workout, interpretError(result.Error)
}
