package controller

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func (c *Controller) loadExerciseRoutes() {
	g := c.baseRouterGroup.Group("/exercises")

	g.GET("/", c.allExercises)
	g.GET("/:id", c.exercise)
	g.GET("/:id/edit", c.exerciseEdit)
	g.GET("/new", c.exerciseNew)
	g.POST("/", c.createExercise)
	g.POST("/:id", c.updateExercise)
	g.POST("/:id/delete", c.deleteExercise)
}

func (c *Controller) allExercises(ctx *gin.Context) {
	exercises, err := c.db.GetAllExercises()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}
	ctx.HTML(
		http.StatusOK,
		"views/exercises.html",
		gin.H{
			"Exercises": exercises,
		},
	)
}

func (c *Controller) exercise(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		redirectToExercises(ctx)
		return
	}

	ex, ok := c.getExerciseAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	ctx.HTML(
		http.StatusOK,
		"views/exercise.html",
		gin.H{
			"Exercise":            ex,
			"RenderedDescription": template.HTML(ex.Description),
		},
	)
}

func (c *Controller) getExerciseAndHandleErrors(id int, ctx *gin.Context) (*data.Exercise, bool) {
	ex, err := c.db.GetExercise(id)
	if err != nil {
		var notFoundErr *data.RecordNotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return nil, false
		}
		panic(err)
	}
	return ex, true
}

func (c *Controller) exerciseEdit(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		redirectToExercises(ctx)
		return
	}

	ex, ok := c.getExerciseAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	viewData := c.getExerciseEditViewData(ex, ctx.Request.URL.Path)

	ctx.HTML(
		http.StatusOK,
		"views/exercise_edit.html",
		viewData,
	)
}

type exerciseEditViewData struct {
	Exercise    *data.Exercise
	Muscles     []data.Muscle
	Equipment   []data.Equipment
	Categories  []data.ExerciseCategory
	RequestPath string
}

func (c *Controller) getExerciseEditViewData(ex *data.Exercise, route string) exerciseEditViewData {
	muscles, err := c.db.GetAllMuscles()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}

	equipment, err := c.db.GetAllEquipment()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}

	categories, err := c.db.GetAllCategories()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}

	return exerciseEditViewData{
		Exercise:    ex,
		Muscles:     muscles,
		Equipment:   equipment,
		Categories:  categories,
		RequestPath: route,
	}
}

func (c *Controller) exerciseNew(ctx *gin.Context) {
	ex := &data.Exercise{}
	viewData := c.getExerciseEditViewData(ex, ctx.Request.URL.Path)
	ctx.HTML(
		http.StatusOK,
		"views/exercise_new.html",
		viewData,
	)
}

func redirectToExercises(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/exercises")
}

func (c *Controller) updateExercise(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		redirectToExercises(ctx)
		return
	}

	ex, ok := c.getExerciseAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	updateExerciseFromForm(ctx, ex)
	c.saveExercise(ex)
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/exercises/%d", id))
}

func updateExerciseFromForm(ctx *gin.Context, ex *data.Exercise) {
	ex.Name = ctx.PostForm("name")
	ex.Description = ctx.PostForm("description")
	ex.SetMuscles(ctx.PostFormArray("muscles"))
	ex.SetSecondaryMuscles(ctx.PostFormArray("sec-muscles"))
	ex.SetEquipment(ctx.PostFormArray("equipment"))
	ex.SetCategory(ctx.PostForm("category"))
}

func (c *Controller) saveExercise(ex *data.Exercise) {
	err := c.db.SaveExercise(ex)
	if err != nil {
		panic(err)
	}
}

func (c *Controller) createExercise(ctx *gin.Context) {
	ex := data.NewExercise()
	updateExerciseFromForm(ctx, ex)
	c.saveExercise(ex)
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/exercises/%d", ex.Id))
}

func (c *Controller) deleteExercise(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		redirectToExercises(ctx)
		return
	}

	ex, ok := c.getExerciseAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	err = c.db.DeleteExercise(ex)
	if err != nil {
		panic(err)
	}
}
