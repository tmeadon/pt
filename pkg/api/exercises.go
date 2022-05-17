package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/db"
)

func (r routes) addExerciseEndpoints(rg *gin.RouterGroup) {
	exercises := rg.Group("/exercises")
	exercises.GET("/", getAllExercises)
	exercises.GET("/:id", getExerciseById)
}

func getAllExercises(ctx *gin.Context) {
	exercises, err := db.ListAllExercises()
	if err != nil {
		panic(err)
	}
	ctx.JSON(200, newListResponse(exercises))
}

func getExerciseById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "id param is not a valid int"})
	}

	exercise, errs := db.GetExerciseById(id)

	for _, e := range errs {
		logger.Error(e.Error())
	}

	ctx.JSON(200, newSingleResponse(exercise))
}