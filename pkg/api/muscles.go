package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/db"
)

func (r routes) addMuscleEndpoints(rg *gin.RouterGroup) {
	muscles := rg.Group("/muscles")
	muscles.GET("/", getAllMuscles)
	muscles.GET("/:id", getMuscle)
}

func getAllMuscles(ctx *gin.Context) {
	muscles, err := db.GetAllMuscles()
	if err != nil {
		panic(err)
	}
	ctx.JSON(200, newResponse(muscles))
}

func getMuscle(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "id param is not a valid int"})
	}

	response := db.GetMuscle(id)
	if response.Status == db.Empty {
		ctx.JSON(http.StatusNotFound, response.Data)
	}
	if response.Status == db.Fail {
		panic(response.Error)
	}

	ctx.JSON(http.StatusOK, newResponse(response.Data))
}
