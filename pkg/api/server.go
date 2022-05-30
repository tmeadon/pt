package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/tmeadon/pt/pkg/data"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Static struct {
	relativePath string
	filePath     string
}

type apiResponse struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  interface{} `json:"results"`
}

var router *gin.Engine
var logger *zap.Logger
var db *data.DB

type routes struct {
	router *gin.Engine
}

func Start(dbPath string, statics []Static) error {
	db = data.InitDatabase(dbPath)
	logger = initLogger()

	root := routes{
		router: initRouter(),
	}

	v1Routes := root.router.Group("api/v1")
	serveEndpoints(root, v1Routes)
	serveStatics(statics)
	return root.router.Run()
}

func initLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	l.Info("initialised logger")
	return l
}

func initRouter() (r *gin.Engine) {
	r = gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	return
}

func serveStatics(statics []Static) {
	for _, s := range statics {
		router.Static(s.relativePath, s.filePath)
	}
}

func serveEndpoints(r routes, rg *gin.RouterGroup) {
	r.addMuscleEndpoints(rg)
	r.addExerciseEndpoints(rg)
	r.addEquipmentEndpoints(rg)
	r.addCategoriesEndpoints(rg)
}

func newResponse[T any](results []T) apiResponse {
	return apiResponse{
		Count:   len(results),
		Results: results,
	}
}

func handleDBError(err error, ctx *gin.Context) {
	if errors.Is(err, &data.RecordNotFoundError{}) {
		ctx.Status(404)
		return
	}
	panic(err)
}

func parseIDParam(ctx *gin.Context) (int, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "id param is not a valid int"})
		return 0, err
	}
	return id, nil
}
