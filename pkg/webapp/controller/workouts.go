package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func loadWorkoutRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/workouts")

	g.GET("/", allWorkouts)
	// g.GET("/:id/edit", workoutEdit)
	g.POST("/", createWorkout)
	// g.POST("/:id", updateWorkout)
	// g.POST("/:id/delete", deleteWorkout)
}

func allWorkouts(ctx *gin.Context) {
	workouts, err := db.GetAllWorkouts()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}

	users, err := db.GetAllUsers()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}

	ctx.HTML(
		http.StatusOK,
		"views/workouts.html",
		gin.H{
			"Workouts": workouts,
			"Users":    users,
		},
	)
}

// func workoutEdit(ctx *gin.Context) {
// 	id, err := parseIDParam(ctx)
// 	if err != nil {
// 		ctx.Redirect(http.StatusFound, "/workouts")
// 		return
// 	}

// 	u, ok := getWorkoutAndHandleErrors(id, ctx)
// 	if !ok {
// 		return
// 	}

// 	ctx.HTML(
// 		http.StatusOK,
// 		"views/workout_edit.html",
// 		u,
// 	)
// }

// func getWorkoutAndHandleErrors(id int, ctx *gin.Context) (*data.Workout, bool) {
// 	u, err := db.GetWorkout(id)
// 	if err != nil {
// 		var notFoundErr *data.RecordNotFoundError
// 		if errors.As(err, &notFoundErr) {
// 			ctx.AbortWithStatus(http.StatusNotFound)
// 			return nil, false
// 		}
// 		panic(err)
// 	}
// 	return u, true
// }

func createWorkout(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.PostForm("userid"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	w := data.NewWorkout(userID)
	w.Name = ctx.PostForm("name")
	saveWorkout(w)
	ctx.Redirect(http.StatusFound, "/workouts")
}

func saveWorkout(w *data.Workout) {
	err := db.SaveWorkout(w)
	if err != nil {
		panic(err)
	}
}

// func updateWorkout(ctx *gin.Context) {
// 	id, err := parseIDParam(ctx)
// 	if err != nil {
// 		ctx.Redirect(http.StatusFound, "/workouts")
// 		return
// 	}

// 	u, ok := getWorkoutAndHandleErrors(id, ctx)
// 	if !ok {
// 		return
// 	}

// 	u.Name = ctx.PostForm("name")
// 	saveWorkout(u)
// 	ctx.Redirect(http.StatusFound, "/workouts")
// }

// func deleteWorkout(ctx *gin.Context) {
// 	id, err := parseIDParam(ctx)
// 	if err != nil {
// 		ctx.Redirect(http.StatusFound, "/workouts")
// 		return
// 	}

// 	u, ok := getWorkoutAndHandleErrors(id, ctx)
// 	if !ok {
// 		return
// 	}

// 	err = db.DeleteWorkout(u)
// 	if err != nil {
// 		panic(err)
// 	}
// }
