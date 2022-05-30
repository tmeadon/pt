package data

func (db DB) GetAllExercises() ([]Exercise, error) {
	var exercises []Exercise
	result := db.gorm.Preload("Category").Preload("Muscles").Preload("SecondaryMuscles").Preload("Equipment").Find(&exercises)
	return exercises, interpretError(result.Error)
}

func (db DB) GetExercise(id int) (Exercise, error) {
	var exercise Exercise
	result := db.gorm.Preload("Category").Preload("Muscles").Preload("SecondaryMuscles").Preload("Equipment").Find(&exercise)
	return exercise, interpretError(result.Error)
}

func (db DB) InsertExercise(exercise *Exercise) (*Exercise, error) {
	result := db.gorm.Create(exercise)
	return exercise, interpretError(result.Error)
}
