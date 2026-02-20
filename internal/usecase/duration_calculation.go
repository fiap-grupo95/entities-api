//go:build legacy
// +build legacy

package usecase

import (
	"errors"
	"math"
	"time"

	"github.com/rs/zerolog/log"
)

// Mova esse mÃ©todo para internal\usecase\service_order_usecase.go

func CalculateExecutionDurationInHours(startedExecutionDate, finalExecutionDate *time.Time) (float64, error) {
	if startedExecutionDate == nil || finalExecutionDate == nil {
		log.Error().Msg("Dates cannot be ni")
		return 0, errors.New("Dates cannot be nil")
	}

	// Check if the final date is earlier than the start date
	if finalExecutionDate.Before(*startedExecutionDate) {
		log.Error().Msg("The final date cannot be earlier than the start date")
		return 0, errors.New("The final date cannot be earlier than the start date")
	}

	// Calculate the duration in hours
	duration := finalExecutionDate.Sub(*startedExecutionDate)
	return round(duration.Hours(), 2), nil
}

func round(value float64, precision int) float64 {
	return math.Round(value*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
}
