package main

import (
	"context"
	"fmt"
	"log"

	"github.com/rwrrioe/geomap/backend/pkg/database"
	"github.com/rwrrioe/geomap/backend/pkg/repository"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := database.DbConnect()
	if err != nil {
		log.Fatal(err)
	}
	err = database.DbMigrate(db)
	if err != nil {
		log.Fatal(err)
	}
	dbRepo := repository.NewProblemRepo(db)
	ans, _ := dbRepo.GetCommonProblemsByDistrict(ctx, 3072807)
	for _, a := range ans {
		fmt.Println(a)
	}
}
