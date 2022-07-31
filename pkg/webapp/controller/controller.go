package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

type Controller struct {
	db              *data.DB
	baseRouterGroup *gin.RouterGroup
}

func NewController(db *data.DB, rg *gin.RouterGroup) *Controller {
	c := Controller{
		db:              db,
		baseRouterGroup: rg,
	}

	c.loadRoutes()
	return &c
}

func (c *Controller) loadRoutes() {
	c.loadExerciseRoutes()
	c.loadHomepageRoutes()
	c.loadUserRoutes()
	c.loadWorkoutRoutes()
}
