package dbseed

import (
	"github.com/tmeadon/pt/pkg/data"
	"github.com/tmeadon/pt/pkg/wger"
)

func fromWgerCategory(wc *wger.ExerciseCategory) *data.ExerciseCategory {
	return &data.ExerciseCategory{
		Id:   wc.Id,
		Name: wc.Name,
	}
}

func fromWgerMuscles(wgerMuscles *[]wger.Muscle) (muscles []data.Muscle) {
	for _, wm := range *wgerMuscles {
		muscles = append(muscles, *fromWgerMuscle(&wm))
	}
	return
}

func fromWgerMuscle(wm *wger.Muscle) *data.Muscle {
	return &data.Muscle{
		Id:         wm.Id,
		Name:       wm.Name,
		SimpleName: wm.SimpleName,
		IsFront:    wm.IsFront,
	}
}

func fromWgerEquipment(wgerEquipment *[]wger.Equipment) (equipment []data.Equipment) {
	for _, we := range *wgerEquipment {
		equipment = append(equipment, *fromWgerEquipmentItem(&we))
	}
	return equipment
}

func fromWgerEquipmentItem(we *wger.Equipment) *data.Equipment {
	return &data.Equipment{
		Id:   we.Id,
		Name: we.Name,
	}
}

func fromWgerExercise(we *wger.Exercise, base *wger.ExerciseBase) *data.Exercise {
	return &data.Exercise{
		Id:               we.Id,
		Name:             we.Name,
		Description:      we.Description,
		Category:         *fromWgerCategory(&base.Category),
		Muscles:          fromWgerMuscles(&base.Muscles),
		SecondaryMuscles: fromWgerMuscles(&base.MusclesSecondary),
		Equipment:        fromWgerEquipment(&base.Equipment),
	}
}
