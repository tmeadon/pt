package dbseed

import (
	"github.com/tmeadon/pt/pkg/db"
	"github.com/tmeadon/pt/pkg/models"
	"github.com/tmeadon/pt/pkg/wger"
)

func SeedFromWger() error {
	err := seedMusclesFromWger()
	if err != nil {
		return err
	}

	err = seedEquipmentFromWger()
	if err != nil {
		return err
	}

	err = seedExerciseCategoriesFromWger()
	if err != nil {
		return err
	}

	err = seedExercisesFromWger()
	if err != nil {
		return err
	}

	return nil
}

func seedMusclesFromWger() error {
	muscles, err := wger.GetAllMuscles()
	if err != nil {
		return err
	}

	for _, wm := range muscles {
		muscle := models.Muscle{
			Id:         wm.Id,
			Name:       wm.Name,
			SimpleName: wm.SimpleName,
			IsFront:    wm.IsFront,
		}
		err = db.InsertMuscle(muscle, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedEquipmentFromWger() error {
	equipment, err := wger.GetAllEquipment()
	if err != nil {
		return err
	}

	for _, we := range equipment {
		e := models.Equipment{
			Id:   we.Id,
			Name: we.Name,
		}
		err = db.InsertEquipment(e, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedExerciseCategoriesFromWger() error {
	cats, err := wger.GetAllExerciseCategories()
	if err != nil {
		return err
	}

	for _, wc := range cats {
		c := models.ExerciseCategory{
			Id:   wc.Id,
			Name: wc.Name,
		}
		err = db.InsertCategory(c, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedExercisesFromWger() error {
	exerciseBases, err := wger.GetAllExerciseBases()
	if err != nil {
		return err
	}

	for _, base := range exerciseBases {
		for _, ex := range filterEnglishExercisesOnly(&base.Exercises) {
			e := models.Exercise{
				Id:               ex.Id,
				Name:             ex.Name,
				Description:      ex.Description,
				Category:         fromWgerCategory(&base.Category),
				Muscles:          fromWgerMuscles(&base.Muscles),
				MusclesSecondary: fromWgerMuscles(&base.MusclesSecondary),
				Equipment:        fromWgerEquipment(&base.Equipment),
			}
			err = db.InsertExercise(e, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
