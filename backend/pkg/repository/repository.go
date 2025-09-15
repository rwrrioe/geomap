package repository

import (
	"context"
	"log"

	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/ybru-tech/georm"
	"gorm.io/gorm"
)

type ProblemStat struct {
	Type  string
	Count int
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
		Type:        p.Type,
		Coordinates: p.Geom,
	}
}

type ProblemRepository interface {
	GetById(ctx context.Context, id int) (*ProblemDTO, error)
	ListByDistrict(ctx context.Context, districtID int) ([]*ProblemDTO, error)
	GetCommonProblemsByDistrict(ctx context.Context, id int) ([]ProblemStat, error)
}

type ProblemRepo struct {
	Db *gorm.DB
}

func NewProblemRepo(db *gorm.DB) *ProblemRepo {
	return &ProblemRepo{Db: db}
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

func (p *ProblemRepo) GetCommonProblemsByDistrict(ctx context.Context, id int) ([]ProblemStat, error) {
	var stats []ProblemStat
	result := p.Db.Raw(
		`SELECT type, 
		COUNT(type) AS count
		FROM problems
		WHERE district_id = ?
		GROUP BY type
		ORDER BY COUNT(type) DESC
		`, id).Scan(&stats)

	if result.Error != nil {
		return nil, result.Error
	}
	return stats, nil
}
