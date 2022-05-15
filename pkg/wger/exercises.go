package wger

import "encoding/json"

type Exercise struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    int    `json:"language"`
}

type ExerciseBase struct {
	Id               int              `json:"id"`
	Category         ExerciseCategory `json:"category"`
	Muscles          []Muscle         `json:"muscles"`
	MusclesSecondary []Muscle         `json:"muscles_secondary"`
	Equipment        []Equipment      `json:"equipment"`
	Exercises        []Exercise       `json:"exercises"`
}

func GetAllExerciseBases() ([]ExerciseBase, error) {
	bases := make([]ExerciseBase, 0)

	url := baseUrl + "/exercisebaseinfo"

	for {
		data, err := get(url)
		if err != nil {
			return nil, err
		}

		dataJson, err := json.Marshal(data.Results)
		if err != nil {
			return nil, err
		}

		baseBatch := make([]ExerciseBase, 0)
		jsonErr := json.Unmarshal(dataJson, &baseBatch)
		if jsonErr != nil {
			return nil, jsonErr
		}

		bases = append(bases, baseBatch...)

		if data.Next == "" {
			break
		}

		url = data.Next
	}

	return bases, nil
}
