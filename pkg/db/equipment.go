package db

import (
	"database/sql"
	"fmt"

	"github.com/tmeadon/pt/pkg/models"
)

func InsertEquipment(newEquipment models.Equipment, keepID bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(getEquipmentInsertStatement(keepID))
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = execEquipmentInsert(stmt, &newEquipment, keepID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func getEquipmentInsertStatement(keepID bool) string {
	if keepID {
		return fmt.Sprintf("INSERT INTO %s (id, name) VALUES (?, ?)", tables.EquipmentTable)
	}
	return fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", tables.EquipmentTable)
}

func execEquipmentInsert(stmt *sql.Stmt, newEquip *models.Equipment, keepID bool) error {
	if keepID {
		_, err := stmt.Exec(newEquip.Id, newEquip.Name)
		return err
	}
	_, err := stmt.Exec(newEquip.Name)
	return err
}

func GetAllEquipment() ([]models.Equipment, error) {
	equipment := make([]models.Equipment, 0)

	query := "select id, name from %s"
	rows, err := queryRows(fmt.Sprintf(query, tables.EquipmentTable))

	if rows == nil {
		return equipment, err
	}

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
