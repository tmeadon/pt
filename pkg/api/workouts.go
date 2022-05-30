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
	workouts.GET("/", getAllWorkouts)
	workouts.GET("/:id", getWorkout)
	workouts.POST("/", newWorkout)
}

func getAllWorkouts(ctx *gin.Context) {
	workouts, err := db.GetAllWorkouts()
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, newResponse(workouts))
}

func getWorkout(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		return
	}

	workout, err := db.GetWorkout(id)
	if err != nil {
		handleDBError(err, ctx)
	}

	ctx.JSON(http.StatusOK, newResponse([]data.Workout{workout}))
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
