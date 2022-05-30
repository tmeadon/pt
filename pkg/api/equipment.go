package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func (r routes) addEquipmentEndpoints(rg *gin.RouterGroup) {
	equipment := rg.Group("/equipment")
	equipment.GET("/", getAllEquipments)
	equipment.GET("/:id", getEquipment)
}

func getAllEquipments(ctx *gin.Context) {
	equipment, err := db.GetAllEquipment()
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(200, newResponse(equipment))
}

func getEquipment(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		return
	}

	equipment, err := db.GetEquipment(id)
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, newResponse([]data.Equipment{equipment}))
}
