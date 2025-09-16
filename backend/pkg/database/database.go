package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/rwrrioe/geomap/backend/pkg/repository"
	"github.com/twpayne/go-geom"
	"github.com/ybru-tech/georm"
	"google.golang.org/genai"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DistrictResponse struct {
	Type      string                 `json:"type"`
	Districts []entities.DistrictDTO `json:"features"`
}

type ProblemsResponse struct {
	Problems []TMPAiResponse `json:"problems"`
}

type TMPAiResponse struct {
	ProblemID   int    `json:"problem_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Geom        string `json:"geom"`
	Status      string `json:"status"`
	Importance  int    `json:"importance"`
}

func newProblemsResponse() *ProblemsResponse {
	return &ProblemsResponse{}
}

func DbConnect() (repository.ProblemRepository, error) {
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
	repo := repository.NewProblemRepo(db)
	return repo, err
}

func DbMigrate(r repository.ProblemRepository) error {
	r.GetDb().Db.AutoMigrate(&entities.District{}, &entities.Problem{})
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

func GenerateProblems(ctx context.Context, db *gorm.DB) error {
	fmt.Println("generating func")
	problemsResponse := newProblemsResponse()
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"problems": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"problem_id":  {Type: genai.TypeNumber},
							"geom":        {Type: genai.TypeString},
							"name":        {Type: genai.TypeString},
							"description": {Type: genai.TypeString},
							"importance":  {Type: genai.TypeNumber},
							"type":        {Type: genai.TypeString},
							"status":      {Type: genai.TypeString},
						},
						Required: []string{"problem_id", "geom", "name", "description", "importance", "type", "status"},
					},
				},
			},
			Required: []string{"problems"},
		},
	}
	prompt := `
Тебе надо сгенерировать 50-70 мелких или средних проблем для проекта карты проблем города Алматы как тестовые данные. 
Важность по шкале от 1 до 10. Статус: created, processing, solved. 
4 типа:  ЖКХ,  Дороги и транспорт,  Гос.сервис, Прочее. 
ProblemID — это первичный ключ и он не должен повторяться.

⚠️ ВАЖНО: координаты должны быть строго в формате POINT(lon lat), 
где сначала указывается долгота (пример: 76.945), потом широта (пример: 43.255). 
Пример: "geom": "POINT(76.945 43.255)".

Пиши строго соблюдая конфиг, без лишних комментариев, 
чтобы всё потом хорошо напарсилось в БД. Колонка geom — это geometry(Point,4326).

Распредели проблемы по всем районам Алматы, чтобы можно было выявить паттерны.
`
	fmt.Println("start generating")
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to generate AI response:%w", err)
	}

	fmt.Println(result, "raw gemini text")
	fmt.Println("end generating, unmarshalling")
	if err := json.Unmarshal([]byte(result.Text()), &problemsResponse); err != nil {
		return fmt.Errorf("failed to unmarshal AI response:%w", err)
	}

	fmt.Println("start problems unparsing")
	for i, p := range problemsResponse.Problems {
		coords := strings.TrimPrefix(strings.TrimSuffix(p.Geom, ")"), "POINT(")
		parts := strings.Split(coords, " ")
		lon, _ := strconv.ParseFloat(parts[0], 64)
		lat, _ := strconv.ParseFloat(parts[1], 64)

		if lon < 60 {
			lon, lat = lat, lon
		}

		point := geom.NewPointFlat(geom.XY, []float64{lon, lat})
		problem := entities.Problem{
			ProblemID:   p.ProblemID,
			Name:        p.Name,
			Description: p.Description,
			Status:      p.Status,
			Importance:  p.Importance,
			Geom:        georm.New(point),
		}
		db.Create(&problem)
		fmt.Println("written", i)
	}
	return nil

}
