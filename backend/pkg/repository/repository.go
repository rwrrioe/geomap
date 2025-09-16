package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
	"github.com/ybru-tech/georm"
	"gorm.io/gorm"
)

type FindDistrictResponse struct {
	District_ID   int    `gorm:"column:district_id"`
	District_name string `gorm:"column:district_name"`
}

type ProblemStatByDistrict struct {
	TypeID       int     `gorm:"column:type_id"`
	TypeName     string  `gorm:"column:type"`
	ProblemCount int     `gorm:"column:prb_count"`
	SolvedCount  int     `gorm:"column:solved_count"`
	ImpAvg       float64 `gorm:"column:avg_imp"`
}

type ProblemStatByType struct {
	DistrictID   int     `gorm:"column:district_id"`
	DistrictName string  `gorm:"column:district_name"`
	ProblemCount int     `gorm:"column:prb_count"`
	StatusCount  int     `gorm:"column:solved_count"`
	ImpAvg       float64 `gorm:"column:avg_imp"`
}

type ProblemStatByCity struct {
	ProblemCount int     `gorm:"column:prb_count"`
	StatusCount  int     `gorm:"column:solved_count"`
	ImpAvg       float64 `gorm:"column:avg_imp"`
}

type ProblemDTO struct {
	ProblemID   int
	Name        string
	Description string
	Type        string
	ImageURL    string //TODO
	Coordinates georm.Point
	Status      string
}

func newProblemDTO(p entities.Problem) *ProblemDTO {
	return &ProblemDTO{
		ProblemID:   p.ProblemID,
		Status:      p.Status,
		Name:        p.Name,
		Description: p.Description,
		Coordinates: p.Geom,
	}
}

type ProblemRepository interface {
	GetById(ctx context.Context, id int) (*ProblemDTO, error)
	ListByDistrict(ctx context.Context, id int) ([]*ProblemDTO, error)
	GetAnalysisByDistrict(ctx context.Context, id int) ([]ProblemStatByDistrict, error)
	GetAnalysisByType(ctx context.Context, id int) ([]ProblemStatByType, error)
	GetAnalysisByCity(ctx context.Context) (ProblemStatByCity, error)
	FindDistrict(ctx context.Context, point geom.Point) (FindDistrictResponse, error)
	GetDb() *ProblemRepo
}

type ProblemRepo struct {
	Db *gorm.DB
}

func NewProblemRepo(db *gorm.DB) *ProblemRepo {
	return &ProblemRepo{Db: db}
}

func (p *ProblemRepo) GetDb() *ProblemRepo {
	return &ProblemRepo{
		Db: p.Db,
	}
}

func (p *ProblemRepo) GetById(ctx context.Context, id int) (*ProblemDTO, error) {
	var problem entities.Problem

	result := p.Db.WithContext(ctx).Preload("District").First(&problem, "problem_id=?", id)
	if result.Error != nil {
		log.Fatal("not found")
	}
	dto := newProblemDTO(problem)
	return dto, nil
}

func (p *ProblemRepo) ListByDistrict(ctx context.Context, id int) ([]*ProblemDTO, error) {
	var problems []entities.Problem

	result := p.Db.WithContext(ctx).Preload("District").Where("district_id=?", id).Find(&problems)
	if result.Error != nil {
		return nil, result.Error
	}

	dtos := make([]*ProblemDTO, 0, len(problems))
	for _, prob := range problems {
		dtos = append(dtos, newProblemDTO(prob))
	}

	return dtos, nil
}

func (p *ProblemRepo) GetAnalysisByDistrict(ctx context.Context, id int) ([]ProblemStatByDistrict, error) {
	var stats []ProblemStatByDistrict
	result := p.Db.Raw(
		`SELECT
                p.type_id,
				type,
                count(problem_id) as prb_count,
                count(problem_id) filter (where status = 'solved') as solved_count,
                avg(importance)::numeric(10,2) as avg_imp
                FROM problems p
				JOIN problem_types USING(type_id)
				 WHERE district_id = ?
                GROUP BY type_id, type`, id).Scan(&stats)

	if result.Error != nil {
		return nil, result.Error
	}
	return stats, nil
}

func (p *ProblemRepo) GetAnalysisByType(ctx context.Context, id int) ([]ProblemStatByType, error) {
	var statByType []ProblemStatByType

	result := p.Db.Raw(
		`SELECT 
		p.district_id,
		d.name_ru as district_name,
		COUNT(problem_id) AS prb_count,
		AVG(p.importance)::numeric(10,2) as avg_imp,
		COUNT(*) FILTER (WHERE p.status = 'solved') as solved_count
		FROM problems p
		JOIN districts d USING(district_id)
		WHERE p.type_id = ?
		GROUP BY p.district_id, d.name_ru
		ORDER BY prb_count DESC
		`, id).Scan(&statByType)

	if result.Error != nil {
		return nil, result.Error
	}

	return statByType, nil
}

func (p *ProblemRepo) GetAnalysisByCity(ctx context.Context) (ProblemStatByCity, error) {
	var cityStat ProblemStatByCity
	result := p.Db.Raw(
		`
		SELECT
		COUNT(problem_id) AS prb_count,
		AVG(importance)::numeric(10,2) as avg_imp,
		COUNT(*) FILTER( WHERE status = 'solved') as solved_count
		FROM problems
		`).Scan(&cityStat)

	if result.Error != nil {
		return ProblemStatByCity{}, result.Error
	}

	return cityStat, nil
}

func (p *ProblemRepo) FindDistrict(ctx context.Context, point geom.Point) (FindDistrictResponse, error) {
	var district FindDistrictResponse

	log.Println("wkt marshalling")
	pointWKT, err := wkt.Marshal(&point)
	if err != nil {
		return FindDistrictResponse{}, err
	}
	fmt.Println(pointWKT)
	log.Println("marshalling ended, query started")
	result := p.Db.Raw(
		`
		SELECT
		district_id,
		name_ru AS district_name
		FROM districts
		WHERE ST_Contains(
		    geom, 
			ST_SetSRID(ST_GeomFromText(?), 4326)
	)
		`, pointWKT).Scan(&district)
	log.Println("query ended")

	log.Println("first check")
	if result.Error != nil {
		return FindDistrictResponse{}, fmt.Errorf("db query failed: %w", result.Error)
	}

	log.Println("second check")
	if result.RowsAffected == 0 {
		return FindDistrictResponse{}, fmt.Errorf("no district found for point %s", pointWKT)
	}
	log.Println("return")
	return district, nil
}

func (p *ProblemRepo) AddProblem(ctx context.Context, problem entities.ProblemResponseDTO) error {
	result := p.Db.Create(&problem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
