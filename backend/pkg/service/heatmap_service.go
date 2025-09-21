package service

import (
	"context"

	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/rwrrioe/geomap/backend/pkg/repository"
)

type HeatMapService struct {
	repo repository.ProblemRepository
}

func NewHeatMapService(repo repository.ProblemRepository) *HeatMapService {
	return &HeatMapService{
		repo: repo,
	}
}

func (h *HeatMapService) BuildHeatMap(ctx context.Context) error {
	points, err := h.repo.ListProblems(ctx)
	if err != nil {
		return err
	}
	heatPoints := make([]entities.HeatPoint, 0, len(*points))

	for _, p := range *points {
		coords := p.Geom.Geom.Coords()
		lon := coords[0]
		lat := coords[1]

		heatPoint := entities.HeatPoint{
			Category: p.TypeID,
			Point: entities.Point{
				Lon:        lon,
				Lat:        lat,
				Importance: p.Importance,
			},
		}

		heatPoints = append(heatPoints, heatPoint)
	}

	heatMap := entities.HeatMap{
		Max:        len(heatPoints),
		HeatPoints: heatPoints,
	}

	err = h.repo.CacheHeatMap(ctx, &heatMap)
	if err != nil {
		return err
	}

	return nil
}

func (h *HeatMapService) GetHeatMap(ctx context.Context) (*entities.CachedHeatMap, error) {
	heatmap, err := h.repo.GetHeatMap(ctx)
	if err != nil {
		err := h.BuildHeatMap(ctx)
		if err != nil {
			return nil, err
		}

		heatmap, err := h.GetHeatMap(ctx)
		if err != nil {
			return nil, err
		}

		return heatmap, nil
	}

	return heatmap, nil
}
