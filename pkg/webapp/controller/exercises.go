package controller

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func loadExerciseRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/exercises")

	g.GET("/", allExercises)
	g.GET("/:id", exercise)
}

func allExercises(ctx *gin.Context) {
	exercises, err := db.GetAllExercises()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}
	ctx.HTML(
		http.StatusOK,
		"exercises.html",
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
		"exercise.html",
		gin.H{
			"Exercise":            ex,
			"RenderedDescription": template.HTML(ex.Description),
		},
	)

}

func redirectToExercises(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "/exercises")
}
