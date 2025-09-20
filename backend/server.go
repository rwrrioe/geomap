package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/geomap/backend/pkg/database"
	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/rwrrioe/geomap/backend/pkg/handlers"
	"github.com/rwrrioe/geomap/backend/pkg/service"
)

type HTTPServer struct {
	*handlers.HTTPHandlers
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) InitServerDefault() error {
	dbRepo, err := database.DbConnect()
	if err != nil {
		return err
	}

	AIService := *service.NewAIPredictService(dbRepo)
	HeatMapService := *service.NewHeatMapService(dbRepo)
	ProblemService := *service.NewProblemService(dbRepo)

	handlers := &handlers.HTTPHandlers{
		User:           &entities.User{Role: "Guest"},
		AIService:      &AIService,
		HeatMapService: &HeatMapService,
		ProblemService: &ProblemService,
	}

	engine := gin.Default()
	s.HTTPHandlers = handlers

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	engine.GET("/heatmap", s.HTTPHandlers.GetHeatmap)
	engine.POST("/heatmap", s.CreateBreefPredicts)
	engine.GET("/heatmap/analysis/district/:districtID", handlers.GetDistrictPrediction)
	engine.GET("/heatmap/analysis/type/:typeID", handlers.GetTypePrediction)
	engine.GET("/heatmap/analysis/city/:cityID", handlers.GetPredictByCity)
	engine.Run(":8080")
	return nil
}
