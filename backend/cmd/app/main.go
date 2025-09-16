package main

import (
	"context"
	"fmt"

	"github.com/rwrrioe/geomap/backend/pkg/database"
	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/rwrrioe/geomap/backend/pkg/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	db, err := database.DbConnect()
	if err != nil {
		fmt.Println(err.Error())
	}
	err = database.DbMigrate(db)
	if err != nil {
		fmt.Println(err.Error())
	}

	req := entities.CreateProblemRequest{
		ProblemID:   72,
		ProblemName: "Test",
		Description: "test",
		TypeID:      3,
		Lat:         76.9457,
		Lon:         43.2389,
	}

	service := service.NewProblemService(db)

	service.NewProblem(ctx, req)
}
