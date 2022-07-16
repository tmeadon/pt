package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseIDParam(ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	return id, err
}
