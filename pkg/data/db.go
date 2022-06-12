package data

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	gorm *gorm.DB
}

func InitDatabase(dbPath string) *DB {
	db := openDatabase(dbPath)
	enableForeignKeys(db)
	applyMigrations(db)
	return &DB{gorm: db}
}

func openDatabase(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to open database - %s", err))
	}
	return db
}

func enableForeignKeys(db *gorm.DB) {
	if res := db.Exec("PRAGMA foreign_keys = ON", nil); res.Error != nil {
		panic(fmt.Errorf("could not enable foreign keys - %s", res.Error))
	}
}

func applyMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&Muscle{}, &Equipment{}, &ExerciseCategory{}, &Exercise{}, &ExerciseHistory{}, &ExerciseSet{}, &User{}, &Workout{})
	if err != nil {
		panic(fmt.Errorf("failed to apply migration - %s", err))
	}
}

func interpretError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &RecordNotFoundError{err: err}
	}
	if err.Error() == "FOREIGN KEY constraint failed" {
		return &ForeignKeyError{err: err}
	}
	return &InternalError{err: err}
}
