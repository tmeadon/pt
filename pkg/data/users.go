package data

import "time"

type User struct {
	Base
	Name     string `json:"name"`
	Username string `json:"username"`
}

func NewUser() *User {
	now := time.Now().UTC()
	return &User{
		Base: Base{
			Created:      now,
			LastModified: now,
		},
	}
}

func (db DB) GetAllUsers() ([]User, error) {
	var users []User
	result := db.gorm.Where("is_deleted = false").Find(&users)
	return users, interpretError(result.Error)
}

func (db DB) GetUser(id int) (*User, error) {
	var user User
	result := db.gorm.Where("is_deleted = false").First(&user, id)
	return &user, interpretError(result.Error)
}

func (db DB) SaveUser(user *User) error {
	user.LastModified = time.Now().UTC()
	result := db.gorm.Save(user)
	return interpretError(result.Error)
}

func (db DB) DeleteUser(user *User) error {
	user.IsDeleted = true
	err := db.SaveUser(user)
	return err
}
