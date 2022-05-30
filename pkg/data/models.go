package data

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
