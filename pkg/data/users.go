package data

import "time"

func (db DB) GetAllUsers() ([]User, error) {
	var users []User
	result := db.gorm.Find(&users)
	return users, interpretError(result.Error)
}

func (db DB) GetUser(id int) (*User, error) {
	var user User
	result := db.gorm.First(&user, id)
	return &user, interpretError(result.Error)
}

func (db DB) InsertUser(user *User) error {
	user.Created = time.Now().UTC()
	user.LastModified = time.Now().UTC()
	result := db.gorm.Create(user)
	return interpretError(result.Error)
}
