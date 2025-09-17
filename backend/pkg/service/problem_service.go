package service

// TODO ADD ListProblems by User, DeleteProblem, edit problem

import (
	"context"

	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/rwrrioe/geomap/backend/pkg/repository"
	"github.com/twpayne/go-geom"
	"github.com/ybru-tech/georm"
)

type ProblemService struct {
	repo repository.ProblemRepository
}

func NewProblemService(repo repository.ProblemRepository) *ProblemService {
	return &ProblemService{
		repo: repo,
	}
}

func (p *ProblemService) NewProblem(ctx context.Context, req entities.CreateProblemRequest) error {
	point := geom.NewPointFlat(geom.XY, []float64{req.Lon, req.Lat})

	district, err := p.repo.GetDb().FindDistrict(ctx, *point)
	if err != nil {
		return err
	}

	problem := entities.Problem{
		DistrictID:  district.District_ID,
		Geom:        georm.New(point),
		Name:        req.ProblemName,
		Description: req.Description,
		Status:      "created",
	}

	err = p.repo.AddProblem(ctx, problem)
	if err != nil {
		return err
	}

	return nil
}
