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

type ExerciseHistory struct {
	Base
	Time       time.Time     `json:"time"`
	ExerciseId int           `json:"-"`
	Exercise   Exercise      `json:"exercise"`
	UserId     int           `json:"-"`
	User       User          `json:"user"`
	Sets       []ExerciseSet `json:"sets"`
	WorkoutId  int           `json:"-"`
}

type ExerciseSet struct {
	Base
	WeightKG          int `json:"weight_kg"`
	Reps              int `json:"reps"`
	ExerciseHistoryId int `json:"-"`
}
