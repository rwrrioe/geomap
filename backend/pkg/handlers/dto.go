package handlers

import (
	"fmt"
	"time"

	"github.com/rwrrioe/geomap/backend/pkg/entities"
)

type errDTO struct {
	Message string
	Time    time.Time
}

func newErrDTO(err error, time time.Time) errDTO {
	return errDTO{
		Message: err.Error(),
		Time:    time,
	}
}

func validateCityID(id int) error {
	if id != -1 {
		return fmt.Errorf("invalid city id")
	}
	return nil
}

func validateProblemReq(req *entities.CreateProblemRequest) error {
	if req.ProblemName == "" {
		return fmt.Errorf("empty problem name")
	}
	return nil
}
