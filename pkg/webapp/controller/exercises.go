package controller

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

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

	ex, err := db.GetExercise(id)
	if err != nil {
		var notFoundErr *data.RecordNotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		panic(err)
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

	viewData := getExerciseViewData(ex)

	ctx.HTML(
		http.StatusOK,
		"views/exercise_edit.html",
		viewData,
	)
}

func exerciseNew(ctx *gin.Context) {
	ex := &data.Exercise{}
	viewData := getExerciseViewData(ex)
	ctx.HTML(
		http.StatusOK,
		"views/exercise_new.html",
		viewData,
	)
}

func redirectToExercises(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/exercises")
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

	ex.Name = ctx.PostForm("name")
	ex.Description = ctx.PostForm("description")

	ex.Muscles = make([]data.Muscle, 0)
	for _, m := range ctx.PostFormArray("muscles") {
		id, _ := strconv.Atoi(m)
		ex.Muscles = append(ex.Muscles, data.Muscle{Base: data.Base{Id: id}})
	}

	ex.SecondaryMuscles = make([]data.Muscle, 0)
	for _, m := range ctx.PostFormArray("sec-muscles") {
		id, _ := strconv.Atoi(m)
		ex.SecondaryMuscles = append(ex.SecondaryMuscles, data.Muscle{Base: data.Base{Id: id}})
	}

	ex.Equipment = make([]data.Equipment, 0)
	for _, m := range ctx.PostFormArray("equipment") {
		id, _ := strconv.Atoi(m)
		ex.Equipment = append(ex.Equipment, data.Equipment{Base: data.Base{Id: id}})
	}

	catId, _ := strconv.Atoi(ctx.PostForm("category"))
	ex.Category = data.ExerciseCategory{Base: data.Base{Id: catId}}

	err = db.UpdateExercise(ex)
	if err != nil {
		panic(err)
	}

	ctx.Request.ParseForm()
	fmt.Println(ctx.Request.Form)

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/exercises/%d", id))
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

func createExercise(ctx *gin.Context) {
	ex := &data.Exercise{
		Name:             ctx.PostForm("name"),
		Description:      ctx.PostForm("description"),
		Muscles:          make([]data.Muscle, 0),
		SecondaryMuscles: make([]data.Muscle, 0),
		Equipment:        make([]data.Equipment, 0),
	}

	for _, m := range ctx.PostFormArray("muscles") {
		id, _ := strconv.Atoi(m)
		ex.Muscles = append(ex.Muscles, data.Muscle{Base: data.Base{Id: id}})
	}

	for _, m := range ctx.PostFormArray("sec-muscles") {
		id, _ := strconv.Atoi(m)
		ex.SecondaryMuscles = append(ex.SecondaryMuscles, data.Muscle{Base: data.Base{Id: id}})
	}

	for _, m := range ctx.PostFormArray("equipment") {
		id, _ := strconv.Atoi(m)
		ex.Equipment = append(ex.Equipment, data.Equipment{Base: data.Base{Id: id}})
	}

	catId, _ := strconv.Atoi(ctx.PostForm("category"))
	ex.Category = data.ExerciseCategory{Base: data.Base{Id: catId}}

	err := db.UpdateExercise(ex)
	if err != nil {
		panic(err)
	}

	ctx.Redirect(http.StatusFound, fmt.Sprintf("/exercises/%d", ex.Id))
}

type exerciseViewData struct {
	Exercise   *data.Exercise
	Muscles    []data.Muscle
	Equipment  []data.Equipment
	Categories []data.ExerciseCategory
}

func getExerciseViewData(ex *data.Exercise) exerciseViewData {
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

	return exerciseViewData{
		Exercise:   ex,
		Muscles:    muscles,
		Equipment:  equipment,
		Categories: categories,
	}
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
