package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func loadUserRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/users")

	g.GET("/", allUsers)
	g.GET("/:id/edit", userEdit)
	g.GET("/new", userNew)
	g.POST("/", createUser)
	g.POST("/:id", updateUser)
	g.POST("/:id/delete", deleteUser)
}

func allUsers(ctx *gin.Context) {
	users, err := db.GetAllUsers()
	if err != nil && !errors.Is(err, &data.RecordNotFoundError{}) {
		panic(err)
	}
	ctx.HTML(
		http.StatusOK,
		"views/users.html",
		gin.H{
			"Users": users,
		},
	)
}

func userEdit(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/users")
		return
	}

	u, ok := getUserAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	ctx.HTML(
		http.StatusOK,
		"views/user_edit.html",
		u,
	)
}

func getUserAndHandleErrors(id int, ctx *gin.Context) (*data.User, bool) {
	u, err := db.GetUser(id)
	if err != nil {
		var notFoundErr *data.RecordNotFoundError
		if errors.As(err, &notFoundErr) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return nil, false
		}
		panic(err)
	}
	return u, true
}

func userNew(ctx *gin.Context) {
	u := &data.User{}
	ctx.HTML(
		http.StatusOK,
		"views/user_new.html",
		u,
	)
}

func createUser(ctx *gin.Context) {
	u := data.NewUser()
	u.Username = ctx.PostForm("username")
	u.Name = ctx.PostForm("name")
	saveUser(u)
	ctx.Redirect(http.StatusFound, "/users")
}

func saveUser(u *data.User) {
	err := db.SaveUser(u)
	if err != nil {
		panic(err)
	}
}

func updateUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/users")
		return
	}

	u, ok := getUserAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	u.Name = ctx.PostForm("name")
	saveUser(u)
	ctx.Redirect(http.StatusFound, "/users")
}

func deleteUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/users")
		return
	}

	u, ok := getUserAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	err = db.DeleteUser(u)
	if err != nil {
		panic(err)
	}
}
