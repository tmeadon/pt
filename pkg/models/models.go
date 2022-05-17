package models

type Muscle struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	SimpleName string `json:"simple_name"`
	IsFront    bool   `json:"is_front"`
}

type Equipment struct {
	Id   int
	Name string
}

type ExerciseCategory struct {
	Id   int
	Name string
}

type Exercise struct {
	Id               int
	Name             string
	Description      string
	Category         ExerciseCategory
	Muscles          []Muscle
	MusclesSecondary []Muscle
	Equipment        []Equipment
}
