package data

func (db DB) GetAllEquipment() ([]Equipment, error) {
	var equipments []Equipment
	result := db.gorm.Find(&equipments)
	return equipments, interpretError(result.Error)
}

func (db DB) GetEquipment(id int) (Equipment, error) {
	var equipment Equipment
	result := db.gorm.First(&equipment, id)
	return equipment, interpretError(result.Error)
}

func (db DB) InsertEquipment(equipment *Equipment) (*Equipment, error) {
	result := db.gorm.Create(equipment)
	return equipment, interpretError(result.Error)
}
