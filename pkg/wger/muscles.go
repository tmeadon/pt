package wger

import (
	"encoding/json"
)

type Muscle struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	SimpleName string `json:"name_en"`
	IsFront    bool   `json:"is_front"`
}

func GetAllMuscles() ([]Muscle, error) {
	muscles := make([]Muscle, 0)

	data, err := get(baseUrl + "/muscle")
	if err != nil {
		return nil, err
	}

	dataJson, err := json.Marshal(data.Results)
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(dataJson, &muscles)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return muscles, nil
}
