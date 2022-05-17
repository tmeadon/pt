package db

import (
	"database/sql"
	"fmt"

	"github.com/tmeadon/pt/pkg/models"
)

func GetCategoryById(id int) (models.ExerciseCategory, error) {
	stmt, err := db.Prepare(fmt.Sprintf("SELECT id, name from %s WHERE id = ?", tables.ExerciseCategoriesTable))
	if err != nil {
		return models.ExerciseCategory{}, err
	}

	cat := models.ExerciseCategory{}

	sqlErr := stmt.QueryRow(id).Scan(&cat.Id, &cat.Name)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return models.ExerciseCategory{}, nil
		}
		return models.ExerciseCategory{}, nil
	}

	return cat, nil
}

func InsertCategory(newCategory models.ExerciseCategory, keepID bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(getCategoryInsertStatement(keepID))
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = execCategoryInsert(stmt, &newCategory, keepID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func getCategoryInsertStatement(keepID bool) string {
	if keepID {
		return fmt.Sprintf("INSERT INTO %s (id, name) VALUES (?, ?)", tables.ExerciseCategoriesTable)
	}
	return fmt.Sprintf("INSERT INTO %s (name) VALUES (?)", tables.ExerciseCategoriesTable)
}

func execCategoryInsert(stmt *sql.Stmt, newCategory *models.ExerciseCategory, keepID bool) error {
	if keepID {
		_, err := stmt.Exec(newCategory.Id, newCategory.Name)
		return err
	}
	_, err := stmt.Exec(newCategory.Name)
	return err
}
