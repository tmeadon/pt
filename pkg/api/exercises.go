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

	history := rg.Group("/exercisehistory")
	history.POST("/", addExerciseHistory)
	history.GET("/:id", getExerciseHistory)
	history.PUT("/:id", updateExerciseHistory)
}

func getAllExercises(ctx *gin.Context) {
	exercises, err := db.GetAllExercises()
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(200, newResponse(exercises))
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

	ctx.JSON(http.StatusOK, newResponse([]data.Exercise{*exercise}))
}

func getExerciseHistory(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		return
	}

	history, err := db.GetExerciseHistory(id)
	if err != nil {
		handleDBError(err, ctx)
	}

	ctx.JSON(http.StatusOK, newResponse([]data.ExerciseHistory{*history}))
}

func addExerciseHistory(ctx *gin.Context) {
	body := exerciseHistoryRequest{}
	if err := validateBody(&body, ctx); err != nil {
		return
	}

	if v := validateExerciseHistoryRequest(&body, ctx); !v {
		return
	}

	history := body.ToModel()

	err := db.SaveExerciseHistory(history)
	if err != nil {
		handleDBError(err, ctx)
	}

	ctx.JSON(http.StatusCreated, newResponse([]data.ExerciseHistory{*history}))
}

func updateExerciseHistory(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		return
	}

	body := exerciseHistoryRequest{}
	if err := validateBody(&body, ctx); err != nil {
		return
	}

	_, err = db.GetExerciseHistory(id)
	if err != nil {
		handleDBError(err, ctx)
	}

	if !validateExerciseHistoryRequest(&body, ctx) {
		return
	}

	history := body.ToModel()
	history.Id = id

	err = db.SaveExerciseHistory(history)
	if err != nil {
		handleDBError(err, ctx)
	}

	ctx.JSON(http.StatusOK, newResponse([]data.ExerciseHistory{*history}))
}
