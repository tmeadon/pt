package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func (r routes) addExerciseEndpoints(rg *gin.RouterGroup) {
	exercises := rg.Group("/exercises")
	exercises.GET("/", getAllExercises)
	exercises.GET("/:id", getExercise)
}

func getAllExercises(ctx *gin.Context) {
	Exercises, err := db.GetAllExercises()
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(200, newResponse(Exercises))
}

func getExercise(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		return
	}

	exercise, err := db.GetExercise(id)
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, newResponse([]data.Exercise{exercise}))
}
