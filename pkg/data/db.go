package data

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	gorm *gorm.DB
}

func InitDatabase(dbPath string) *DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to open database")
	}

	err = db.AutoMigrate(&Muscle{}, &Equipment{}, &ExerciseCategory{}, &Exercise{})
	if err != nil {
		panic("failed to apply migration")
	}

	return &DB{gorm: db}
}

func interpretError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &RecordNotFoundError{}
	}
	return &InternalError{err: err}
}
