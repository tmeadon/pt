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

	return nil
}

func seedMusclesFromWger() error {
	err := db.RecreateMusclesTable()
	if err != nil {
		return err
	}

	muscles, err := wger.GetAllMuscles()
	if err != nil {
		return err
	}

	for i := 0; i < len(muscles); i++ {
		muscle := models.Muscle{
			Id:         muscles[i].Id,
			Name:       muscles[i].Name,
			SimpleName: muscles[i].SimpleName,
			IsFront:    muscles[i].IsFront,
		}
		err = db.InsertMuscle(muscle, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func seedEquipmentFromWger() error {
	err := db.RecreateEquipmentTable()
	if err != nil {
		return err
	}

	equipment, err := wger.GetAllEquipment()
	if err != nil {
		return err
	}

	for i := 0; i < len(equipment); i++ {
		e := models.Equipment{
			Id:   equipment[i].Id,
			Name: equipment[i].Name,
		}
		err = db.InsertEquipment(e, true)
		if err != nil {
			return err
		}
	}

	return nil
}
