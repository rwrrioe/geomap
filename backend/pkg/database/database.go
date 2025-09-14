package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rwrrioe/geomap/pkg/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DistrictResponse struct {
	Type      string                 `json:"type"`
	Districts []entities.DistrictDTO `json:"features"`
}

func DbConnect() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsnParam := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	dsn := dsnParam
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

func DbMigrate(db *gorm.DB) error {
	db.AutoMigrate(&entities.District{}, &entities.Problem{})
	return nil
}

func ParseDistrict(db *gorm.DB) error {
	var districtsResponse DistrictResponse
	file, err := os.Open("almaty.json")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&districtsResponse); err != nil {
		return fmt.Errorf("error %w", err)
	}

	for i, d := range districtsResponse.Districts {
		district := entities.MapToDistinct(d)
		db.Create(&district)
		fmt.Println("parsed to the db, #", i)
	}

	return nil
}
