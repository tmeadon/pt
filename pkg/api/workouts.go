package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

type newWorkoutRequest struct {
	Name   string `json:"name" binding:"required"`
	UserId int    `json:"user_id" binding:"required"`
}

func (r routes) addWorkoutEndpoints(rg *gin.RouterGroup) {
	workouts := rg.Group("/workouts")
	workouts.POST("/", newWorkout)
}

func newWorkout(ctx *gin.Context) {
	body := newWorkoutRequest{}
	if err := validateBody(&body, ctx); err != nil {
		return
	}
	workout, err := db.InsertWorkout(&data.Workout{Name: body.Name, UserId: body.UserId})
	if err != nil {
		handleDBError(err, ctx)
	}
	ctx.JSON(http.StatusAccepted, &workout)
}
