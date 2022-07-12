package api

import (
	"fmt"
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

	ctx.JSON(http.StatusOK, newResponse([]data.User{*user}))
}

func newUser(ctx *gin.Context) {
	body := newWorkoutRequest{}
	if err := validateBody(&body, ctx); err != nil {
		return
	}

	// check user exists
	if _, err := db.GetUser(body.UserId); err != nil {
		m := fmt.Sprintf("User %d does not exist", body.UserId)
		logger.Info(m)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": m})
		return
	}

	workout := data.Workout{Name: body.Name, UserId: body.UserId}
	err := db.InsertWorkout(&workout)
	if err != nil {
		handleDBError(err, ctx)
		return
	}
	ctx.JSON(http.StatusCreated, newResponse([]data.Workout{workout}))
}
