package data

import "time"

func (db DB) GetAllCategories() ([]ExerciseCategory, error) {
	var categories []ExerciseCategory
	result := db.gorm.Find(&categories)
	return categories, interpretError(result.Error)
}

func (db DB) GetCategory(id int) (ExerciseCategory, error) {
	var category ExerciseCategory
	result := db.gorm.First(&category, id)
	return category, interpretError(result.Error)
}

func (db DB) InsertCategory(category *ExerciseCategory) (*ExerciseCategory, error) {
	category.Created = time.Now().UTC()
	category.LastModified = time.Now().UTC()
	result := db.gorm.Create(category)
	return category, interpretError(result.Error)
}
