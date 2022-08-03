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

func (c *Controller) loadWorkoutRoutes() {
	g := c.baseRouterGroup.Group("/workouts")

	g.GET("/", c.allWorkouts)
	g.GET("/:id", c.workout)
	g.POST("/", c.createWorkout)
	g.POST("/:id/delete", c.deleteWorkout)
	g.POST("/:id/exercise", c.addWorkoutExercise)
	g.POST("/:id/exercise/:exId/delete", c.deleteWorkoutExercise)
	g.POST("/:id/exercise/:exId/sets", c.addWorkoutSet)
	g.POST("/:id/sets/:setId", c.updateWorkoutSet)
	g.POST("/:id/sets/:setId/delete", c.deleteWorkoutSet)
}

func (c *Controller) allWorkouts(ctx *gin.Context) {
	workouts, err := c.db.GetAllWorkouts()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}

	users, err := c.db.GetAllUsers()
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

func (c *Controller) workout(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	w, ok := c.getWorkoutAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	ex, err := c.db.GetAllExercises()
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

func (c *Controller) getWorkoutAndHandleErrors(id int, ctx *gin.Context) (*data.Workout, bool) {
	u, err := c.db.GetWorkout(id)
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

func (c *Controller) getSetAndHandleError(id int, ctx *gin.Context) (*data.ExerciseSet, bool) {
	s, err := c.db.GetSet(id)
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

func (c *Controller) getExerciseHistoryAndHandleError(id int, ctx *gin.Context) (*data.ExerciseHistory, bool) {
	e, err := c.db.GetExerciseHistory(id)
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

func (c *Controller) createWorkout(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.PostForm("userid"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	w := data.NewWorkout(userID)
	c.saveWorkout(w)
	ctx.Redirect(http.StatusFound, "/workouts")
}

func (c *Controller) saveWorkout(w *data.Workout) {
	err := c.db.SaveWorkout(w)
	if err != nil {
		panic(err)
	}
}

func (c *Controller) deleteWorkout(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	w, ok := c.getWorkoutAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	err = c.db.DeleteWorkout(w)
	if err != nil {
		panic(err)
	}
}

func (c *Controller) addWorkoutExercise(ctx *gin.Context) {
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

	w, ok := c.getWorkoutAndHandleErrors(workoutId, ctx)
	if !ok {
		return
	}

	w.ExerciseInstances = append(w.ExerciseInstances, *data.NewExerciseHistory(w.UserId, w.Id, exId))

	if err = c.db.SaveWorkout(w); err != nil {
		panic(err)
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}

func (c *Controller) addWorkoutSet(ctx *gin.Context) {
	workoutId, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	exId, err := parseIntParam(ctx, "exId")
	if err != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
	}

	ex, ok := c.getExerciseHistoryAndHandleError(exId, ctx)
	if !ok {
		return
	}

	weight, weightErr := strconv.ParseFloat(ctx.PostForm("weight"), 32)
	reps, repsErr := strconv.Atoi(ctx.PostForm("reps"))

	if weightErr != nil || repsErr != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
	}

	ex.AddSet(data.NewExerciseSet(float32(weight), reps))

	if err = c.db.SaveExerciseHistory(ex); err != nil {
		panic(err)
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}

func (c *Controller) updateWorkoutSet(ctx *gin.Context) {
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

	set, ok := c.getSetAndHandleError(setId, ctx)
	if !ok {
		return
	}

	weight, weightErr := strconv.ParseFloat(ctx.PostForm("weight"), 32)
	reps, repsErr := strconv.Atoi(ctx.PostForm("reps"))

	if weightErr != nil || repsErr != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
	}

	set.WeightKG = float32(weight)
	set.Reps = reps

	if err := c.db.SaveSet(set); err != nil {
		panic(err)
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}

func (c *Controller) deleteWorkoutSet(ctx *gin.Context) {
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

	set, ok := c.getSetAndHandleError(setId, ctx)
	if !ok {
		return
	}

	if err := c.db.DeleteSet(set); err != nil {
		panic(err)
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}

func (c *Controller) deleteWorkoutExercise(ctx *gin.Context) {
	workoutId, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/workouts")
		return
	}

	exId, err := parseIntParam(ctx, "exId")
	if err != nil {
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
	}

	ex, ok := c.getExerciseHistoryAndHandleError(exId, ctx)
	if !ok {
		return
	}

	c.db.DeleteExerciseHistory(ex)
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/workouts/%d", workoutId))
}
