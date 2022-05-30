package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func (r routes) addUserEndpoints(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("/", getAllUsers)
	users.GET("/:id", getUser)
}

func getAllUsers(ctx *gin.Context) {
	users, err := db.GetAllUsers()
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, newResponse(users))
}

func getUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		return
	}

	user, err := db.GetUser(id)
	if err != nil {
		handleDBError(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, newResponse([]data.User{user}))
}
