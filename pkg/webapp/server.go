package webapp

import (
	"html/template"
	"strings"

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

	s.gin.SetFuncMap(template.FuncMap{
		"joinMuscleNames": joinMuscleNames,
	})

	s.gin.LoadHTMLGlob("web/templates/*")
	s.gin.Static("/css", "./web/static/css")
	s.gin.Static("/img", "./web/static/img")
	s.gin.Static("/scss", "./web/static/scss")
	s.gin.Static("/vendor", "./web/static/vendor")
	s.gin.Static("/js", "./web/static/js")
	s.gin.StaticFile("/favicon.ico", "./web/img/favicon.ico")

	// r.GET("/", func(ctx *gin.Context) {
	// 	exercises, err := db.GetAllExercises()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	ctx.HTML(
	// 		http.StatusOK,
	// 		"index.html",
	// 		gin.H{
	// 			"Exercises": exercises,
	// 		},
	// 	)
	// })
	s.gin.Run()
	return nil
}

func joinMuscleNames(muscles []data.Muscle) string {
	names := make([]string, 0)
	for _, m := range muscles {
		if m.SimpleName != "" {
			names = append(names, m.SimpleName)
			continue
		}
		names = append(names, m.Name)
	}
	return strings.Join(names, ", ")
}
