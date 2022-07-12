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

type Exercise struct {
	Base
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Category         ExerciseCategory `json:"category"`
	CategoryId       int              `json:"-"`
	Muscles          []Muscle         `json:"muscles" gorm:"many2many:exercise_muscles"`
	SecondaryMuscles []Muscle         `json:"secondary_muscles" gorm:"many2many:exercise_secondary_muscles"`
	Equipment        []Equipment      `json:"equipment" gorm:"many2many:exercise_equipment"`
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

type User struct {
	Base
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Workout struct {
	Base
	Name               string             `json:"name"`
	UserId             int                `json:"user_id"`
	User               User               `json:"user"`
	ExerciseInstances  []ExerciseHistory  `json:"exercises"`
	ExerciseCategories []ExerciseCategory `json:"exercise_categories" gorm:"many2many:workout_categories"`
}
