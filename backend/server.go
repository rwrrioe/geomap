package server

import (
	"log"

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

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	uploadsPath := "root/uploads"
	log.Println(uploadsPath)
	engine.Static("/uploads", uploadsPath)

	s.HTTPHandlers = handlers

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{""},
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
	engine.GET("/heatmap/districts/:districtID/problems/:problemID", handlers.GetProblem)
	engine.GET("/heatmap/districts/:districtID/problems", handlers.ListProblemsByDistrict)
	engine.POST("/heatmap/districts/:districtID/problems", handlers.CreateProblem)
	engine.Run(":8080")
	return nil
}
