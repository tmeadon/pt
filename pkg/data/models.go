package data

import "time"

type Muscle struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	SimpleName string `json:"simple_name"`
	IsFront    bool   `json:"is_front"`
}

type Equipment struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ExerciseCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Exercise struct {
	Id               int              `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Category         ExerciseCategory `json:"category"`
	CategoryId       int              `json:"-"`
	Muscles          []Muscle         `json:"muscles" gorm:"many2many:exercise_muscles"`
	SecondaryMuscles []Muscle         `json:"secondary_muscles" gorm:"many2many:exercise_secondary_muscles"`
	Equipment        []Equipment      `json:"equipment" gorm:"many2many:exercise_equipment"`
}

type ExerciseInstance struct {
	Id         int           `json:"id"`
	Time       time.Time     `json:"time"`
	ExerciseId int           `json:"-"`
	Exercise   Exercise      `json:"exercise"`
	UserId     int           `json:"-"`
	User       User          `json:"user"`
	Sets       []ExerciseSet `json:"sets"`
	WorkoutId  int           `json:"-"`
}

type ExerciseSet struct {
	Id                 int `json:"id"`
	WeightKG           int `json:"weight_kg"`
	Reps               int `json:"reps"`
	ExerciseInstanceId int `json:"-"`
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type Workout struct {
	Id                 int                `json:"id"`
	CreatedAt          time.Time          `json:"created_at"`
	UserId             int                `json:"-"`
	User               User               `json:"user"`
	ExerciseInstances  []ExerciseInstance `json:"workouts"`
	ExerciseCategories []ExerciseCategory `json:"exercise_categories" gorm:"many2many:workout_categories"`
}
