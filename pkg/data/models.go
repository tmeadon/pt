package data

import "time"

type Base struct {
	Id           int       `json:"id"`
	Created      time.Time `json:"-"`
	LastModified time.Time `json:"-"`
	IsDeleted    bool      `json:"-" gorm:"default:false"`
}

type Muscle struct {
	Base
	Name       string `json:"name"`
	SimpleName string `json:"simple_name"`
	IsFront    bool   `json:"is_front"`
}

type Equipment struct {
	Base
	Name string `json:"name"`
}

type ExerciseCategory struct {
	Base
	Name string `json:"name"`
}
