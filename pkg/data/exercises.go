package data

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

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

func NewExercise() *Exercise {
	now := time.Now().UTC()
	return &Exercise{
		Base: Base{
			Created:      now,
			LastModified: now,
		},
	}
}

func (ex *Exercise) SetMuscles(muscleIDs []string) error {
	ex.Muscles = make([]Muscle, 0)
	for _, id := range muscleIDs {
		i, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		ex.Muscles = append(ex.Muscles, Muscle{Base: Base{Id: i}})
	}
	return nil
}

func (ex *Exercise) SetSecondaryMuscles(muscleIDs []string) error {
	ex.SecondaryMuscles = make([]Muscle, 0)
	for _, id := range muscleIDs {
		i, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		ex.SecondaryMuscles = append(ex.SecondaryMuscles, Muscle{Base: Base{Id: i}})
	}
	return nil
}

func (ex *Exercise) SetEquipment(equipmentIDs []string) error {
	ex.Equipment = make([]Equipment, 0)
	for _, id := range equipmentIDs {
		i, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		ex.Equipment = append(ex.Equipment, Equipment{Base: Base{Id: i}})
	}
	return nil
}

func (ex *Exercise) SetCategory(categoryID string) error {
	id, err := strconv.Atoi(categoryID)
	if err != nil {
		return err
	}
	ex.Category = ExerciseCategory{Base: Base{Id: id}}
	return nil
}

func (db *DB) GetAllExercises() ([]Exercise, error) {
	var exercises []Exercise
	result := db.gorm.Where("is_deleted = false").Preload("Category").Preload("Muscles").Preload("SecondaryMuscles").Preload("Equipment").Find(&exercises)
	return exercises, interpretError(result.Error)
}

func (db *DB) GetExercise(id int) (*Exercise, error) {
	var exercise Exercise
	result := db.gorm.Where("is_deleted = false").Preload("Category").Preload("Muscles").Preload("SecondaryMuscles").Preload("Equipment").First(&exercise, id)
	return &exercise, interpretError(result.Error)
}

func (db *DB) SaveExercise(ex *Exercise) error {
	ex.LastModified = time.Now().UTC()
	err := db.gorm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(ex).Error; err != nil {
			return err
		}

		if err := tx.Model(ex).Association("Muscles").Replace(ex.Muscles); err != nil {
			return err
		}

		if err := tx.Model(ex).Association("SecondaryMuscles").Replace(ex.SecondaryMuscles); err != nil {
			return err
		}

		if err := tx.Model(ex).Association("Equipment").Replace(ex.Equipment); err != nil {
			return err
		}

		return nil
	})
	return interpretError(err)
}

func (db DB) DeleteExercise(exercise *Exercise) error {
	exercise.IsDeleted = true
	result := db.gorm.Save(exercise)
	return interpretError(result.Error)
}
