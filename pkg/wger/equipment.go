package wger

import (
	"encoding/json"
)

type Equipment struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllEquipment() ([]Equipment, error) {
	equipment := make([]Equipment, 0)

	data, err := get(baseUrl + "/equipment")
	if err != nil {
		return nil, err
	}

	dataJson, err := json.Marshal(data.Results)
	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(dataJson, &equipment)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return equipment, nil
}
