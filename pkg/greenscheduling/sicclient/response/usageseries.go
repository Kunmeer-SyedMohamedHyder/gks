package response

import (
	"fmt"
	"time"
)

// UsageSeriesItem represents a single item in the usage series response.
type UsageSeriesItem struct {
	ID            string   `json:"id"`
	Type          string   `json:"type"`
	TimeBucket    string   `json:"timeBucket"`
	CostUsd       *float64 `json:"costUsd"`       // Use pointer to allow for null
	Co2eMetricTon *float64 `json:"co2eMetricTon"` // Use pointer to allow for null
	Kwh           *float64 `json:"kwh"`           // Use pointer to allow for null
}

// GetTimeBucket returns the TimeBucket value.
func (u *UsageSeriesItem) GetTimeBucket() (time.Time, error) {
	parsedTime, err := time.Parse(time.RFC3339, u.TimeBucket)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing time: %w", err)
	}

	return parsedTime, nil
}

// GetCostUsd returns the CostUsd value or zero if it's nil.
func (u *UsageSeriesItem) GetCostUsd() float64 {
	if u.CostUsd != nil {
		return *u.CostUsd
	}
	return 0.0 // or any other default value you'd like
}

// GetCo2eMetricTon returns the Co2eMetricTon value or zero if it's nil.
func (u *UsageSeriesItem) GetCo2eMetricTon() float64 {
	if u.Co2eMetricTon != nil {
		return *u.Co2eMetricTon
	}
	return 0.0 // or any other default value you'd like
}

// GetKwh returns the Kwh value or zero if it's nil.
func (u *UsageSeriesItem) GetKwh() float64 {
	if u.Kwh != nil {
		return *u.Kwh
	}
	return 0.0 // or any other default value you'd like
}

// UsageSeriesResponse represents the response structure for usage series data.
type UsageSeriesResponse struct {
	Items []UsageSeriesItem `json:"items"`
	Count int               `json:"count"`
}
