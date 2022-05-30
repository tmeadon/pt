package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func (r routes) addCategoriesEndpoints(rg *gin.RouterGroup) {
	categories := rg.Group("/categories")
	categories.GET("/", getAllCategories)
	categories.GET("/:id", getCategory)
}

func getAllCategories(ctx *gin.Context) {
	categories, err := db.GetAllCategories()
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(200, newResponse(categories))
}

func getCategory(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		return
	}

	category, err := db.GetCategory(id)
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, newResponse([]data.ExerciseCategory{category}))
}
