package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
)

func (c *Controller) loadUserRoutes() {
	g := c.baseRouterGroup.Group("/users")

	g.GET("/", c.allUsers)
	g.GET("/:id/edit", c.userEdit)
	g.GET("/new", c.userNew)
	g.POST("/", c.createUser)
	g.POST("/:id", c.updateUser)
	g.POST("/:id/delete", c.deleteUser)
}

func (c *Controller) allUsers(ctx *gin.Context) {
	users, err := c.db.GetAllUsers()
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

func (c *Controller) userEdit(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/users")
		return
	}

	u, ok := c.getUserAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	ctx.HTML(
		http.StatusOK,
		"views/user_edit.html",
		u,
	)
}

func (c *Controller) getUserAndHandleErrors(id int, ctx *gin.Context) (*data.User, bool) {
	u, err := c.db.GetUser(id)
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

func (c *Controller) userNew(ctx *gin.Context) {
	u := &data.User{}
	ctx.HTML(
		http.StatusOK,
		"views/user_new.html",
		u,
	)
}

func (c *Controller) createUser(ctx *gin.Context) {
	u := data.NewUser()
	u.Username = ctx.PostForm("username")
	u.Name = ctx.PostForm("name")
	c.saveUser(u)
	ctx.Redirect(http.StatusFound, "/users")
}

func (c *Controller) saveUser(u *data.User) {
	err := c.db.SaveUser(u)
	if err != nil {
		panic(err)
	}
}

func (c *Controller) updateUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/users")
		return
	}

	u, ok := c.getUserAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	u.Name = ctx.PostForm("name")
	c.saveUser(u)
	ctx.Redirect(http.StatusFound, "/users")
}

func (c *Controller) deleteUser(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/users")
		return
	}

	u, ok := c.getUserAndHandleErrors(id, ctx)
	if !ok {
		return
	}

	err = c.db.DeleteUser(u)
	if err != nil {
		panic(err)
	}
}
