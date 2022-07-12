package dbseed

import (
	"github.com/tmeadon/pt/pkg/data"
	"github.com/tmeadon/pt/pkg/wger"
)

func fromWgerCategory(wc *wger.ExerciseCategory) *data.ExerciseCategory {
	ec := &data.ExerciseCategory{
		Name: wc.Name,
	}
	ec.Id = wc.Id
	return ec
}

func fromWgerMuscles(wgerMuscles *[]wger.Muscle) (muscles []data.Muscle) {
	for _, wm := range *wgerMuscles {
		muscles = append(muscles, *fromWgerMuscle(&wm))
	}
	return
}

func fromWgerMuscle(wm *wger.Muscle) *data.Muscle {
	m := &data.Muscle{
		Name:       wm.Name,
		SimpleName: wm.SimpleName,
		IsFront:    wm.IsFront,
	}
	m.Id = wm.Id
	return m
}

func fromWgerEquipment(wgerEquipment *[]wger.Equipment) (equipment []data.Equipment) {
	for _, we := range *wgerEquipment {
		equipment = append(equipment, *fromWgerEquipmentItem(&we))
	}
	return equipment
}

func fromWgerEquipmentItem(we *wger.Equipment) *data.Equipment {
	eq := &data.Equipment{
		Name: we.Name,
	}
	eq.Id = we.Id
	return eq
}

func fromWgerExercise(we *wger.Exercise, base *wger.ExerciseBase) *data.Exercise {
	ex := &data.Exercise{
		Name:             we.Name,
		Description:      we.Description,
		Category:         *fromWgerCategory(&base.Category),
		Muscles:          fromWgerMuscles(&base.Muscles),
		SecondaryMuscles: fromWgerMuscles(&base.MusclesSecondary),
		Equipment:        fromWgerEquipment(&base.Equipment),
	}
	ex.Id = we.Id
	return ex
}
