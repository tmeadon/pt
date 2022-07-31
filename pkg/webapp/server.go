package webapp

import (
	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
	"github.com/tmeadon/pt/pkg/webapp/controller"
)

type Server struct {
	db         *data.DB
	gin        *gin.Engine
	controller *controller.Controller
}

func NewServer(dbPath string) *Server {
	db := data.InitDatabase(dbPath)
	g := initGin()
	c := controller.NewController(db, g.Group("/"))

	return &Server{
		db:         data.InitDatabase(dbPath),
		gin:        g,
		controller: c,
	}
}

func initGin() (g *gin.Engine) {
	g = gin.New()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())
	return
}

func (s *Server) Start() error {
	s.gin.Use(gin.Logger())
	s.gin.Use(gin.Recovery())

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
