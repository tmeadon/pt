package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
	"github.com/tmeadon/pt/pkg/webapp/controller"
)

var server *Server

type Server struct {
	db  *data.DB
	gin *gin.Engine
}

func NewServer(dbPath string) *Server {
	server = &Server{
		db:  data.InitDatabase(dbPath),
		gin: gin.New(),
	}
	return server
}

func (s *Server) Start() error {
	s.gin.Use(gin.Logger())
	s.gin.Use(gin.Recovery())
	controller.LoadRoutes(s.db, s.gin.Group("/"))

	s.gin.SetFuncMap(controller.TemplateFuncs())

	s.gin.LoadHTMLGlob("web/templates/**/*")
	s.gin.Static("/css", "./web/static/css")
	s.gin.Static("/img", "./web/static/img")
	s.gin.Static("/scss", "./web/static/scss")
	s.gin.Static("/vendor", "./web/static/vendor")
	s.gin.Static("/js", "./web/static/js")
	s.gin.StaticFile("/favicon.ico", "./web/img/favicon.ico")

	s.gin.Run()
	return nil
}
