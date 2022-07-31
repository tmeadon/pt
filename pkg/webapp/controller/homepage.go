package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *Controller) loadHomepageRoutes() {
	c.baseRouterGroup.GET("/", c.homepage)
}

func (c *Controller) homepage(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"views/homepage.html",
		nil,
	)
}
