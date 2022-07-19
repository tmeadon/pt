package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func loadHomepageRoutes(rg *gin.RouterGroup) {
	rg.GET("/", homepage)
}

func homepage(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"views/homepage.html",
		nil,
	)
}
