package handlers

import (
	"fmt"
	"time"
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
