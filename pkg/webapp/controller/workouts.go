package controller

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func loadWorkoutRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/workouts")

	g.GET("/", allWorkouts)
	g.GET("/:id", workout)
	g.POST("/", createWorkout)
	g.POST("/:id/delete", deleteWorkout)
	g.POST("/:id/exercise", addWorkoutExercise)
	g.POST("/:id/exercise/:exId/delete", deleteWorkoutExercise)
	g.POST("/:id/exercise/:exId/sets", addWorkoutSet)
	g.POST("/:id/sets/:setId", updateWorkoutSet)
	g.POST("/:id/sets/:setId/delete", deleteWorkoutSet)
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

func workout(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	w, ok := getWorkoutAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	ex, err := db.GetAllExercises()
	if err != nil {
		panic(err)
	}

	// sort exercises by name
	sort.Slice(ex, func(i, j int) bool {
		return ex[i].Name < ex[j].Name
	})

	ctx.HTML(
		http.StatusOK,
		"views/workout.html",
		gin.H{
			"Workout":   w,
			"Exercises": ex,
		},
	)
}

func getWorkoutAndHandleErrors(id int, ctx *gin.Context) (*data.Workout, bool) {
	u, err := db.GetWorkout(id)
	if err != nil {
		var notFoundErr *data.RecordNotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return nil, false
		}
		panic(err)
	}
	return u, true
}

func getSetAndHandleError(id int, ctx *gin.Context) (*data.ExerciseSet, bool) {
	s, err := db.GetSet(id)
	if err != nil {
		var notFoundErr *data.RecordNotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return nil, false
		}
		panic(err)
	}
	return s, true
}

func getExerciseHistoryAndHandleError(id int, ctx *gin.Context) (*data.ExerciseHistory, bool) {
	e, err := db.GetExerciseHistory(id)
	if err != nil {
		var notFoundErr *data.RecordNotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return nil, false
		}
		panic(err)
	}
	return e, true
}

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

func deleteWorkout(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	w, ok := getWorkoutAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	err = db.DeleteWorkout(w)
	if err != nil {
		panic(err)
	}
}

func addWorkoutExercise(ctx *gin.Context) {
	workoutId, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	exId, err := strconv.Atoi(ctx.PostForm("exerciseid"))
	if err != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
		return
	}

	w, ok := getWorkoutAndHandleErrors(workoutId, ctx)
	if !ok {
		return
	}

	w.ExerciseInstances = append(w.ExerciseInstances, *data.NewExerciseHistory(w.UserId, w.Id, exId))

	if err = db.SaveWorkout(w); err != nil {
		panic(err)
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}

func addWorkoutSet(ctx *gin.Context) {
	workoutId, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	exId, err := parseIntParam(ctx, "exId")
	if err != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
	}

	ex, ok := getExerciseHistoryAndHandleError(exId, ctx)
	if !ok {
		return
	}

	weight, weightErr := strconv.Atoi(ctx.PostForm("weight"))
	reps, repsErr := strconv.Atoi(ctx.PostForm("reps"))

	if weightErr != nil || repsErr != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
	}

	ex.AddSet(data.NewExerciseSet(weight, reps))

	if err = db.SaveExerciseHistory(ex); err != nil {
		panic(err)
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}

func updateWorkoutSet(ctx *gin.Context) {
	workoutId, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	setId, err := parseIntParam(ctx, "setId")
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	set, ok := getSetAndHandleError(setId, ctx)
	if !ok {
		return
	}

	weight, weightErr := strconv.Atoi(ctx.PostForm("weight"))
	reps, repsErr := strconv.Atoi(ctx.PostForm("reps"))

	if weightErr != nil || repsErr != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
	}

	set.WeightKG = weight
	set.Reps = reps

	if err := db.SaveSet(set); err != nil {
		panic(err)
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}

func deleteWorkoutSet(ctx *gin.Context) {
	workoutId, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	setId, err := parseIntParam(ctx, "setId")
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	set, ok := getSetAndHandleError(setId, ctx)
	if !ok {
		return
	}

	if err := db.DeleteSet(set); err != nil {
		panic(err)
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}

func deleteWorkoutExercise(ctx *gin.Context) {
	workoutId, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	exId, err := parseIntParam(ctx, "exId")
	if err != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
	}

	ex, ok := getExerciseHistoryAndHandleError(exId, ctx)
	if !ok {
		return
	}

	db.DeleteExerciseHistory(ex)
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}
