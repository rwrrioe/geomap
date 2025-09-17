package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rwrrioe/geomap/backend/pkg/database"
	"github.com/rwrrioe/geomap/backend/pkg/service"
)

func main() {
	r := gin.Default()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	db, err := database.DbConnect()
	if err != nil {
		fmt.Println(err.Error())
	}

	r.POST("/bycity", func(c *gin.Context) {
		aiservice := service.NewAIPredictService(db)

		analysis, err := aiservice.GetAnalysisByCity(ctx)
		if err != nil {
			c.JSON(500, err.Error())
		}
		c.JSON(200, analysis)
	})

	r.Run(":8080")
}
