package service

// TODO ADD ListProblems by User, DeleteProblem, edit problem

import (
	"context"
	"errors"
	"fmt"

	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/rwrrioe/geomap/backend/pkg/repository"
	"github.com/twpayne/go-geom"
	"github.com/ybru-tech/georm"
	"gorm.io/gorm"
)

type ProblemService struct {
	repo repository.ProblemRepository
}

func NewProblemService(repo repository.ProblemRepository) *ProblemService {
	return &ProblemService{
		repo: repo,
	}
}

func (p *ProblemService) NewProblem(ctx context.Context, req entities.CreateProblemForm) error {
	point := geom.NewPointFlat(geom.XY, []float64{req.Lon, req.Lat})

	district, err := p.repo.GetDb().FindDistrict(ctx, *point)
	if err != nil {
		return err
	}

	if ok := p.repo.IsProblemType(ctx, req.TypeID); !ok {
		return fmt.Errorf("invalid problem type")
	}

	problem := entities.Problem{
		DistrictID:  district.District_ID,
		Geom:        georm.New(point),
		Name:        req.ProblemName,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Status:      "created",
		TypeId:      req.TypeID,
	}

	err = p.repo.AddProblem(ctx, problem)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProblemService) GetProblem(ctx context.Context, problemId int) (*repository.ProblemDTO, error) {
	problem, err := p.repo.GetById(ctx, problemId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	return problem, nil
}

func (p *ProblemService) ListProblemsByDistrict(ctx context.Context, districtId int) (*[]repository.ProblemDTO, error) {
	problems, err := p.repo.ListByDistrict(ctx, districtId)
	if err != nil {
		return nil, err
	}

	return problems, nil
}
