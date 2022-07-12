package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func validateBody[T any](body T, ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"message": "Invalid inputs. Please check your inputs",
			},
		)
		return err
	}
	return nil
}

func validateExerciseHistoryRequest(h *exerciseHistoryRequest, ctx *gin.Context) bool {
	if !validateExerciseExists(h.ExerciseId, ctx) {
		return false
	}

	if validateWorkoutExists(h.WorkoutId, ctx) {
		return false
	}

	if validateUserExists(h.UserId, ctx) {
		return false
	}

	return true
}

func validateExerciseExists(exerciseId int, ctx *gin.Context) bool {
	if _, err := db.GetExercise(exerciseId); err != nil {
		if errors.Is(err, &data.RecordNotFoundError{}) {
			ctx.JSON(400, gin.H{"message": fmt.Sprintf("exercise with ID %d not found", exerciseId)})
			return false
		}
	}
	return true
}

func validateWorkoutExists(workoutId int, ctx *gin.Context) bool {
	if _, err := db.GetWorkout(workoutId); err != nil {
		if errors.Is(err, &data.RecordNotFoundError{}) {
			ctx.JSON(400, gin.H{"message": fmt.Sprintf("workout with ID %d not found", workoutId)})
			return false
		}
	}
	return true
}

func validateUserExists(userId int, ctx *gin.Context) bool {
	if _, err := db.GetUser(userId); err != nil {
		if errors.Is(err, &data.RecordNotFoundError{}) {
			ctx.JSON(400, gin.H{"message": fmt.Sprintf("user with ID %d not found", userId)})
			return false
		}
	}
	return true
}
