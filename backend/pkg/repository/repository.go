package repository

//TODO add precached list of all problems
//TODO add redis caching

import (
	"context"
	"errors"
	"fmt"

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
	ProblemID    int         `json:"problem_id"`
	DistrictName string      `json:"district_name"`
	Geom         georm.Point `json:"geom,omitempty"`
	Name         string      `json:"problem_name"`
	Description  string      `json:"problem_desc"`
	ImageURL     string      `gorm:"column:image_url" json:"image_url"`
	Importance   float64     `json:"column:importance"`
	Status       string      `json:"status"`
	TypeID       int         `json:"problem_typeid"`
}

func newProblemDTO(ctx context.Context, repo ProblemRepository, p *entities.Problem) (*ProblemDTO, error) {
	district, err := repo.FindDistrict(ctx, *p.Geom.Geom)
	if err != nil {
		return nil, err
	}

	return &ProblemDTO{
		ProblemID:    p.ProblemID,
		DistrictName: district.District_name,
		Geom:         p.Geom,
		Name:         p.Name,
		Description:  p.Description,
		Importance:   p.Importance,
		TypeID:       p.TypeId,

		Status: p.Status,
	}, nil
}

type ProblemRepository interface {
	GetById(ctx context.Context, id int) (*ProblemDTO, error)
	ListByDistrict(ctx context.Context, id int) (*[]ProblemDTO, error)
	GetAnalysisByDistrict(ctx context.Context, id int) ([]ProblemStatByDistrict, error)
	GetAnalysisByType(ctx context.Context, id int) ([]ProblemStatByType, error)
	GetAnalysisByCity(ctx context.Context) (ProblemStatByCity, error)
	FindDistrict(ctx context.Context, point geom.Point) (FindDistrictResponse, error)
	AddProblem(ctx context.Context, problem entities.Problem) error
	ListProblems(ctx context.Context) (*[]ProblemDTO, error)
	GetAIResponseById(ctx context.Context, id int) (*entities.CachedAnswer, error)
	CacheAIResponse(ctx context.Context, aiResponse *entities.ExtendedAIResponse, requestID int) error
	CacheHeatMap(ctx context.Context, heatmap *entities.HeatMap) error
	GetHeatMap(ctx context.Context) (*entities.CachedHeatMap, error)
	IsDistrict(ctx context.Context, id int) bool
	IsProblemType(ctx context.Context, id int) bool
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

func (p *ProblemRepo) IsDistrict(ctx context.Context, id int) bool {
	result := p.Db.First(&entities.District{}, id)

	return result.Error == nil
}

func (p *ProblemRepo) IsProblemType(ctx context.Context, id int) bool {
	result := p.Db.WithContext(ctx).First(&entities.ProblemType{}, id)

	return result.Error == nil
}

func (p *ProblemRepo) GetAIResponseById(ctx context.Context, id int) (*entities.CachedAnswer, error) {
	var extendedAnswer entities.CachedAnswer

	result := p.Db.Model(&entities.CachedAnswer{}).Find(&entities.CachedAnswer{}, id).Scan(&extendedAnswer)
	if result.Error != nil {
		return nil, result.Error
	}

	return &extendedAnswer, nil
}

func (p *ProblemRepo) CacheAIResponse(ctx context.Context, aiResponse *entities.ExtendedAIResponse, requestID int) error {
	cachedResponse := entities.CachedAnswer{
		ResponseText: aiResponse.AnswerText,
		Status:       aiResponse.Status,
		RequestID:    requestID,
	}

	result := p.Db.Create(&cachedResponse)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *ProblemRepo) CacheHeatMap(ctx context.Context, heatmap *entities.HeatMap) error {
	cachedHeatMap := entities.CachedHeatMap{
		HeatMap: entities.HeatMap{
			Max:        heatmap.Max,
			HeatPoints: heatmap.HeatPoints,
		},
	}

	result := p.Db.Create(&cachedHeatMap)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *ProblemRepo) GetHeatMap(ctx context.Context) (*entities.CachedHeatMap, error) {
	var heatmap entities.CachedHeatMap

	result := p.Db.WithContext(ctx).Scan(&heatmap)
	if result.Error != nil {
		return nil, result.Error
	}

	return &heatmap, nil
}

func (p *ProblemRepo) GetById(ctx context.Context, id int) (*ProblemDTO, error) {
	var problem entities.Problem

	result := p.Db.WithContext(ctx).First(&problem, id)
	if result.Error != nil {
		return nil, result.Error
	}

	dto, err := newProblemDTO(ctx, p, &problem)
	if err != nil {
		return nil, err
	}

	return dto, nil
}

func (p *ProblemRepo) ListByDistrict(ctx context.Context, id int) (*[]ProblemDTO, error) {
	var problems []entities.Problem

	result := p.Db.WithContext(ctx).Preload("District").Where("district_id=?", id).Find(&problems)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	dtos := make([]ProblemDTO, 0, len(problems))
	for _, prob := range problems {
		newDTO, err := newProblemDTO(ctx, p, &prob)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, *newDTO)
	}

	return &dtos, nil
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

	pointWKT, err := wkt.NewEncoder().Encode(&point)
	if err != nil {
		return FindDistrictResponse{}, err
	}
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

	if result.Error != nil {
		return FindDistrictResponse{}, fmt.Errorf("db query failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return FindDistrictResponse{}, fmt.Errorf("no district found for point %s", pointWKT)
	}

	return district, nil
}

func (p *ProblemRepo) AddProblem(ctx context.Context, problem entities.Problem) error {
	result := p.Db.Create(&problem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *ProblemRepo) ListProblems(ctx context.Context) (*[]ProblemDTO, error) {
	var problems []entities.Problem
	result := p.Db.WithContext(ctx).Find(&problems)
	if result.Error != nil {
		return nil, result.Error
	}

	dtos := make([]ProblemDTO, 0, len(problems))

	for _, pr := range problems {
		dto, err := newProblemDTO(ctx, p, &pr)
		if err != nil {
			return nil, err
		}

		dtos = append(dtos, *dto)
	}

	return &dtos, nil
}
