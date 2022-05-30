package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func (r routes) addMuscleEndpoints(rg *gin.RouterGroup) {
	muscles := rg.Group("/muscles")
	muscles.GET("/", getAllMuscles)
	muscles.GET("/:id", getMuscle)
}

func getAllMuscles(ctx *gin.Context) {
	muscles, err := db.GetAllMuscles()
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(200, newResponse(muscles))
}

func getMuscle(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		return
	}

	muscle, err := db.GetMuscle(id)
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, newResponse([]data.Muscle{muscle}))
}
