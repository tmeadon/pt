package api

import (
	"time"

	"github.com/tmeadon/pt/pkg/data"
)

type exerciseHistoryRequest struct {
	Time       time.Time `json:"time" binding:"required"`
	ExerciseId int       `json:"exercise_id" binding:"required"`
	UserId     int       `json:"user_id" binding:"required"`
	WorkoutId  int       `json:"workout_id" binding:"required"`
	Sets       []set     `json:"sets"`
}

func (e *exerciseHistoryRequest) ToModel() *data.ExerciseHistory {
	sets := make([]data.ExerciseSet, 0)
	for _, s := range e.Sets {
		sets = append(sets, *s.ToModel())
	}

	return &data.ExerciseHistory{
		Time:       e.Time,
		ExerciseId: e.ExerciseId,
		UserId:     e.UserId,
		WorkoutId:  e.WorkoutId,
		Sets:       sets,
	}
}

type set struct {
	WeightKG int `json:"weight_kg"`
	Reps     int `json:"reps"`
}

func (s *set) ToModel() *data.ExerciseSet {
	return &data.ExerciseSet{
		WeightKG: s.WeightKG,
		Reps:     s.Reps,
	}
}

type newWorkoutRequest struct {
	Name   string `json:"name" binding:"required"`
	UserId int    `json:"user_id" binding:"required"`
}

type newUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
}

func (ur *newUserRequest) ToModel() *data.User {
	return &data.User{
		Name:     ur.Name,
		Username: ur.Username,
	}
}
