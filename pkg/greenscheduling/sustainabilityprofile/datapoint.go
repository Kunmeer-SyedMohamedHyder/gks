package sustainabilityprofile

import (
	"math"
	"time"
)

// DecayParameters holds the parameters needed for decay calculations.
type DecayParameters struct {
	LatestTime time.Time // Latest time against which decay is calculated
	DecayRate  float64   // Decay rate
}

// NewDecayParameters creates a new instance of DecayParameters.
func NewDecayParameters(latestTime time.Time, decayRate float64) DecayParameters {
	return DecayParameters{
		LatestTime: latestTime,
		DecayRate:  decayRate,
	}
}

// EmissionDataPoint holds data for a single CO₂ emission entry.
type EmissionDataPoint struct {
	Co2  float64   // CO₂ emission value (metric tons)
	Time time.Time // Timestamp of the emission data point
}

// NewEmissionDataPoint creates a new instance of EmissionDataPoint.
func NewEmissionDataPoint(co2 float64, timestamp time.Time) EmissionDataPoint {
	return EmissionDataPoint{
		Co2:  co2,
		Time: timestamp,
	}
}

// CalculateDecayFactor calculates the decay factor for the emission data point
// using the decay parameters.
func (e *EmissionDataPoint) CalculateDecayFactor(params DecayParameters) float64 {
	// Calculate the time difference in hours
	timeDifference := params.LatestTime.Sub(e.Time).Hours()

	// Apply the decay formula: e^(-decayRate*(latestTime-ti)) / (1 + e^(-decayRate*(latestTime-ti)))
	decay := 1.0 // Default decay value in case of invalid time difference
	if timeDifference >= 0 {
		decay = (math.Exp(-params.DecayRate * timeDifference)) / (1 + math.Exp(-params.DecayRate*timeDifference))
	}

	return decay
}
