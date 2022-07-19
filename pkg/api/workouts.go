package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func (r routes) addWorkoutEndpoints(rg *gin.RouterGroup) {
	workouts := rg.Group("/workouts")
	workouts.GET("/", getAllWorkouts)
	workouts.POST("/", newWorkout)
	workouts.GET("/:id", getWorkout)
	workouts.DELETE("/:id", deleteWorkout)
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
		return
	}

	ctx.JSON(http.StatusOK, newResponse([]data.Workout{*workout}))
}

func newWorkout(ctx *gin.Context) {
	body := newWorkoutRequest{}
	if err := validateBody(&body, ctx); err != nil {
		return
	}

	// check user exists
	if _, err := db.GetUser(body.UserId); err != nil {
		m := fmt.Sprintf("User %d does not exist", body.UserId)
		logger.Info(m)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": m})
		return
	}

	workout := data.Workout{Name: body.Name, UserId: body.UserId}
	err := db.SaveWorkout(&workout)
	if err != nil {
		handleDBError(err, ctx)
		return
	}
	ctx.JSON(http.StatusCreated, newResponse([]data.Workout{workout}))
}

func deleteWorkout(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		return
	}

	workout, err := db.GetWorkout(id)
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	if err := db.DeleteWorkout(workout); err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}
