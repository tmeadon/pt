package dbseed

import (
	"github.com/tmeadon/pt/pkg/data"
	"github.com/tmeadon/pt/pkg/wger"
)

var db *data.DB

func SeedFromWger(d *data.DB) error {
	db = d
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
		m := fromWgerMuscle(&wm)
		_, err = db.InsertMuscle(m)
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
		e := fromWgerEquipmentItem(&we)
		_, err = db.InsertEquipment(e)
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
		c := fromWgerCategory(&wc)
		_, err = db.InsertCategory(c)
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
			e := fromWgerExercise(&ex, &base)
			err = db.InsertExercise(e)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
