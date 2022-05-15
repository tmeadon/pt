package dbseed

import (
	"github.com/tmeadon/pt/pkg/models"
	"github.com/tmeadon/pt/pkg/wger"
)

func fromWgerCategory(wgerCat *wger.ExerciseCategory) (cat models.ExerciseCategory) {
	cat.Id = wgerCat.Id
	cat.Name = wgerCat.Name
	return
}

func fromWgerMuscles(wgerMuscles *[]wger.Muscle) (muscles []models.Muscle) {
	for _, wm := range *wgerMuscles {
		muscles = append(muscles, models.Muscle{
			Id:         wm.Id,
			Name:       wm.Name,
			SimpleName: wm.SimpleName,
			IsFront:    wm.IsFront,
		})
	}
	return
}

func fromWgerEquipment(wgerEquipment *[]wger.Equipment) (equipment []models.Equipment) {
	for _, we := range *wgerEquipment {
		equipment = append(equipment, models.Equipment{
			Id:   we.Id,
			Name: we.Name,
		})
	}
	return
}
