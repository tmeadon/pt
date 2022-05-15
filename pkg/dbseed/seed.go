package dbseed

import (
	"github.com/tmeadon/pt/pkg/db"
	"github.com/tmeadon/pt/pkg/models"
	"github.com/tmeadon/pt/pkg/wger"
)

func SeedMusclesFromWger() error {
	err := db.RecreateMusclesTable()
	if err != nil {
		return err
	}

	muscles, err := wger.GetAllMuscles()
	if err != nil {
		return err
	}

	for i := 0; i < len(muscles); i++ {
		err = db.InsertMuscle(models.Muscle{
			Name:       muscles[i].Name,
			SimpleName: muscles[i].SimpleName,
			IsFront:    muscles[i].IsFront,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func SeedEquipmentFromWger() error {
	err := db.RecreateEquipmentTable()
	if err != nil {
		return err
	}

	equipment, err := wger.GetAllEquipment()
	if err != nil {
		return err
	}

	for i := 0; i < len(equipment); i++ {
		err = db.InsertEquipment(models.Equipment{
			Name: equipment[i].Name,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
