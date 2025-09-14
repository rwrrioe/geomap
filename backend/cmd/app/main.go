package main

import (
	"log"

	"github.com/rwrrioe/geomap/pkg/database"
)

func main() {
	db, err := database.DbConnect()
	if err != nil {
		log.Fatal(err)
	}

	if err := database.DbMigrate(db); err != nil {
		log.Fatal(err)
	}

	err = database.ParseDistrict(db)
	if err != nil {
		log.Fatal(err)
	}
}
