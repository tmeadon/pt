package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/db"
)

func (r routes) addMuscleEndpoints(rg *gin.RouterGroup) {
	muscles := rg.Group("/muscles")
	muscles.GET("/", getAllMuscles)
}

func getAllMuscles(ctx *gin.Context) {
	muscles, err := db.GetAllMuscles()
	if err != nil {
		panic(err)
	}
	ctx.JSON(200, newListResponse(muscles))
}
