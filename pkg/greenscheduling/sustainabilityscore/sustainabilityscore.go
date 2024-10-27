package sustainabilityscore

// SustainabilityProfile holds the complete data for sustainability score calculations,
// including emission data points, total CO₂, and total cost values.
type SustainabilityProfile struct {
	Emissions []EmissionDataPoint // Individual CO₂ and time entries
	TotalCo2  float64             // Total CO₂ emissions (metric tons), provided by the user
	TotalCost float64             // Total cost (USD), provided by the user
}

// NewSustainabilityProfile creates a new instance of SustainabilityProfile .
func NewSustainabilityProfile(emissions []EmissionDataPoint, totalCo2, totalCost float64) SustainabilityProfile {
	return SustainabilityProfile{
		Emissions: emissions,
		TotalCo2:  totalCo2,
		TotalCost: totalCost,
	}
}

// CalculateScore calculates the overall sustainability score based on the given weights.
func (data *SustainabilityProfile) CalculateScore(weights SustainabilityWeights) float64 {
	co2WeightedScore := data.calculateCO2WeightedScore(weights.CO2DecayWeight, weights.DecayRate)
	totalCO2WeightedScore := data.calculateTotalCO2WeightedScore(weights.TotalCO2Weight)
	costWeightedScore := data.calculateCostWeightedScore(weights.CostWeight)

	return co2WeightedScore + totalCO2WeightedScore + costWeightedScore
}

// calculateCO2WeightedScore calculates the CO₂ weighted score using the decay parameters.
func (data *SustainabilityProfile) calculateCO2WeightedScore(co2Weight float64, decayRate float64) float64 {
	n := len(data.Emissions)
	if n == 0 {
		return 0
	}

	weightedCO2 := 0.0

	// Get the latest timestamp from the emissions data
	latestTime := data.Emissions[0].Time
	for _, emission := range data.Emissions {
		if emission.Time.After(latestTime) {
			latestTime = emission.Time
		}
	}

	for _, emission := range data.Emissions {
		decayFactor := emission.CalculateDecayFactor(NewDecayParameters(
			latestTime,
			decayRate,
		))
		weightedCO2 += decayFactor * emission.Co2
	}

	return co2Weight * weightedCO2
}

// calculateTotalCO2WeightedScore calculates the total CO₂ weighted score based on the total emissions.
func (data *SustainabilityProfile) calculateTotalCO2WeightedScore(totalCO2Weight float64) float64 {
	return totalCO2Weight * data.TotalCo2
}

// calculateCostWeightedScore calculates the cost weighted score based on the total cost.
func (data *SustainabilityProfile) calculateCostWeightedScore(costWeight float64) float64 {
	return costWeight / (1 + data.TotalCost)
}
