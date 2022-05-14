package wger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var baseUrl string = "https://wger.de/api/v2"

type responseBase struct {
	Count    int                      `json:"count"`
	Next     string                   `json:"next"`
	Previous string                   `json:"previous"`
	Results  []map[string]interface{} `json:"results"`
}

func get(url string) (responseBase, error) {
	fmt.Println("Calling URL: " + url)

	resp, err := http.Get(url)
	if err != nil {
		return responseBase{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseBase{}, nil
	}

	data := responseBase{}

	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		return responseBase{}, jsonErr
	}

	return data, nil
}
