package data

import (
	"time"
)

type ExerciseSet struct {
	Base
	WeightKG          float32 `json:"weight_kg"`
	Reps              int     `json:"reps"`
	ExerciseHistoryId int     `json:"-"`
}

func NewExerciseSet(weightKg float32, reps int) *ExerciseSet {
	now := time.Now().UTC()
	return &ExerciseSet{
		Base: Base{
			Created:      now,
			LastModified: now,
		},
		WeightKG: weightKg,
		Reps:     reps,
	}
}

func (db DB) GetSet(id int) (*ExerciseSet, error) {
	var set ExerciseSet
	result := db.gorm.Where("is_deleted = false").First(&set, id)
	return &set, interpretError(result.Error)
}

func (db DB) SaveSet(set *ExerciseSet) error {
	set.LastModified = time.Now().UTC()
	result := db.gorm.Save(set)
	return interpretError(result.Error)
}

func (db DB) DeleteSet(set *ExerciseSet) error {
	result := db.gorm.Delete(set)
	return interpretError(result.Error)
}
