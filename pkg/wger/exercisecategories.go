package wger

import (
	"encoding/json"
)

type ExerciseCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllExerciseCategories() ([]ExerciseCategory, error) {
	categories := make([]ExerciseCategory, 0)

	data, err := get(baseUrl + "/exercisecategory")
	if err != nil {
		return nil, err
	}

	dataJson, err := json.Marshal(data.Results)
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(dataJson, &categories)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return categories, nil
}
