package controller

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func loadExerciseRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/exercises")

	g.GET("/", allExercises)
	g.GET("/:id", exercise)
	g.GET("/:id/edit", exerciseEdit)
	g.GET("/new", exerciseNew)
	g.POST("/", createExercise)
	g.POST("/:id", updateExercise)
	g.POST("/:id/delete", deleteExercise)
}

func allExercises(ctx *gin.Context) {
	exercises, err := db.GetAllExercises()
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

func exercise(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		redirectToExercises(ctx)
		return
	}

	ex, ok := getExerciseAndHandleErrors(id, ctx)
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

func getExerciseAndHandleErrors(id int, ctx *gin.Context) (*data.Exercise, bool) {
	ex, err := db.GetExercise(id)
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

func exerciseEdit(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		redirectToExercises(ctx)
		return
	}

	ex, ok := getExerciseAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	viewData := getExerciseEditViewData(ex, ctx.Request.URL.Path)

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

func getExerciseEditViewData(ex *data.Exercise, route string) exerciseEditViewData {
	muscles, err := db.GetAllMuscles()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}

	equipment, err := db.GetAllEquipment()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}

	categories, err := db.GetAllCategories()
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

func exerciseNew(ctx *gin.Context) {
	ex := &data.Exercise{}
	viewData := getExerciseEditViewData(ex, ctx.Request.URL.Path)
	ctx.HTML(
		http.StatusOK,
		"views/exercise_new.html",
		viewData,
	)
}

func redirectToExercises(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/exercises")
}

func updateExercise(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		redirectToExercises(ctx)
		return
	}

	ex, ok := getExerciseAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	updateExerciseFromForm(ctx, ex)
	saveExercise(ex)
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

func saveExercise(ex *data.Exercise) {
	err := db.SaveExercise(ex)
	if err != nil {
		panic(err)
	}
}

func createExercise(ctx *gin.Context) {
	ex := data.NewExercise()
	updateExerciseFromForm(ctx, ex)
	saveExercise(ex)
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/exercises/%d", ex.Id))
}

func deleteExercise(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		redirectToExercises(ctx)
		return
	}

	ex, ok := getExerciseAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	err = db.DeleteExercise(ex)
	if err != nil {
		panic(err)
	}
}
