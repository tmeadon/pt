package data

func (db DB) GetAllMuscles() ([]Muscle, error) {
	var muscles []Muscle
	result := db.gorm.Find(&muscles)
	return muscles, interpretError(result.Error)
}

func (db DB) GetMuscle(id int) (Muscle, error) {
	var muscle Muscle
	result := db.gorm.First(&muscle, id)
	return muscle, interpretError(result.Error)
}

func (db DB) InsertMuscle(muscle *Muscle) (*Muscle, error) {
	result := db.gorm.Create(muscle)
	return muscle, interpretError(result.Error)
}
