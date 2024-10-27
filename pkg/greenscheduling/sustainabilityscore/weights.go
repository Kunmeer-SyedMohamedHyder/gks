package sustainabilityscore

// WeightType represents the different types of weights used in sustainability scoring.
type WeightType int

const (
	// Define weight types as constants
	CO2Decay  WeightType = iota // 0
	TotalCO2                    // 1
	Cost                        // 2
	DecayRate                   // 3
)

// SustainabilityWeights holds the weights for different components of the sustainability score.
type SustainabilityWeights struct {
	CO2DecayWeight float64 // Weight for the CO₂ emissions with decay
	TotalCO2Weight float64 // Weight for the total CO₂ emissions
	CostWeight     float64 // Weight for the cost
	DecayRate      float64 // Decay rate for CO₂ emissions
}

// NewSustainabilityWeights creates a new instance of SustainabilityWeights.
// It requires all weight parameters to be provided.
func NewSustainabilityWeights(co2DecayWeight, totalCO2Weight, costWeight, decayRate float64) SustainabilityWeights {
	return SustainabilityWeights{
		CO2DecayWeight: co2DecayWeight,
		TotalCO2Weight: totalCO2Weight,
		CostWeight:     costWeight,
		DecayRate:      decayRate,
	}
}

// GetWeight returns the appropriate weight based on the WeightType.
func (w SustainabilityWeights) GetWeight(weightType WeightType) float64 {
	switch weightType {
	case CO2Decay:
		return w.CO2DecayWeight
	case TotalCO2:
		return w.TotalCO2Weight
	case Cost:
		return w.CostWeight
	case DecayRate:
		return w.DecayRate
	default:
		return 0.0 // Return 0 for an invalid weight type
	}
}
