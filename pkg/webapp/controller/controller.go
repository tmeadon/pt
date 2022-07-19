package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

var db *data.DB

func LoadRoutes(d *data.DB, rg *gin.RouterGroup) {
	db = d
	loadHomepageRoutes(rg)
	loadExerciseRoutes(rg)
	loadUserRoutes(rg)
	loadWorkoutRoutes(rg)
}
