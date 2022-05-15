package db

import (
	"database/sql"
	"fmt"

	"github.com/tmeadon/pt/pkg/models"
)

func InsertEquipment(newEquipment models.Equipment) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", tables.equipmentTable))
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newEquipment.Name)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func GetAllEquipment() ([]models.Equipment, error) {
	equipment := make([]models.Equipment, 0)

	stmt, err := DB.Prepare(fmt.Sprintf("SELECT id, name, full_name, is_front from %s", tables.equipmentTable))
	if err != nil {
		return []models.Equipment{}, err
	}

	rows, err := stmt.Query()
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Equipment{}, nil
		}
		return []models.Equipment{}, err
	}
	defer rows.Close()

	for rows.Next() {
		eq := new(models.Equipment)
		err := rows.Scan(&eq.Id, &eq.Name)
		if err != nil {
			return []models.Equipment{}, err
		}
		equipment = append(equipment, *eq)
	}

	return equipment, nil
}
