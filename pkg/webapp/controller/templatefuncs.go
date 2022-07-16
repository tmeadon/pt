package controller

import (
	"html/template"
	"strings"

	"github.com/tmeadon/pt/pkg/data"
)

func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"ExerciseContainsMuscle":      exerciseContainsMuscle,
		"ExerciseContainsSecMuscle":   exerciseContainsSecMuscle,
		"ExerciseContainsEquipment":   exerciseContainsEquipment,
		"ExerciseHasSecondaryMuscles": exerciseHasSecondaryMuscles,
		"JoinMuscleNames":             joinMuscleNames,
	}
}

func exerciseContainsMuscle(exercise data.Exercise, muscle data.Muscle) bool {
	for _, m := range exercise.Muscles {
		if muscle.Id == m.Id {
			return true
		}
	}
	return false
}

func exerciseContainsSecMuscle(exercise data.Exercise, muscle data.Muscle) bool {
	for _, m := range exercise.SecondaryMuscles {
		if muscle.Id == m.Id {
			return true
		}
	}
	return false
}

func exerciseContainsEquipment(exercise data.Exercise, equipment data.Equipment) bool {
	for _, m := range exercise.Equipment {
		if equipment.Id == m.Id {
			return true
		}
	}
	return false
}

func exerciseHasSecondaryMuscles(exercise data.Exercise) bool {
	return len(exercise.SecondaryMuscles) > 0
}

func joinMuscleNames(muscles []data.Muscle) string {
	names := make([]string, 0)
	for _, m := range muscles {
		if m.SimpleName != "" {
			names = append(names, m.SimpleName)
			continue
		}
		names = append(names, m.Name)
	}
	return strings.Join(names, ", ")
}
