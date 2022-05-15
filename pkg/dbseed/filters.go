package dbseed

import "github.com/tmeadon/pt/pkg/wger"

func filterEnglishExercisesOnly(all *[]wger.Exercise) (filtered []wger.Exercise) {
	for _, e := range *all {
		if e.Language == 2 { // 2 is the id for the english language in the wger api
			filtered = append(filtered, e)
		}
	}
	return
}
