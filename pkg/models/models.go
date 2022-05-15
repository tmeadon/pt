package models

type Muscle struct {
	Id         int
	Name       string
	SimpleName string
	IsFront    bool
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

// type ExerciseBase struct {
// 	Id               int
// 	Category         ExerciseCategory
// 	Muscles          []Muscle
// 	MusclesSecondary []Muscle
// 	Equipment        []Equipment
// 	Exercises        []Exercise
// }
