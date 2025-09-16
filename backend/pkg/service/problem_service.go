package service

import (
	"context"
	"log"

	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/rwrrioe/geomap/backend/pkg/repository"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
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
	point := *geom.NewPointFlat(geom.XY, []float64{req.Lon, req.Lat})
	log.Println("problemservice: start find district")
	district, err := p.repo.GetDb().FindDistrict(ctx, point)
	if err != nil {
		return err
	}

	log.Println("problemservice: end, start converting")
	problem := entities.ProblemResponseDTO{
		ProblemID:   req.ProblemID,
		DistrictID:  district.District_ID,
		Geom:        point,
		Name:        req.ProblemName,
		Description: req.Description,
		Status:      "created",
	}
	localWKT, _ := wkt.Marshal(&point)
	log.Println("local marshalling", localWKT)
	log.Println("problemservice: end, add problem")
	err = p.repo.GetDb().AddProblem(ctx, problem)
	if err != nil {
		return err
	}

	return nil
}
