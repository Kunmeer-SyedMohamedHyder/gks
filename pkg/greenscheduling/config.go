package greenscheduling

// TimeSeriesConfig holds time-related configurations for time series data collection.
type TimeSeriesConfig struct {
	DaysToConsider float64
	SeriesInterval string
}

// SustainabilityWeights holds weighting factors for sustainability metrics.
type SustainabilityWeights struct {
	CO2DecayWeight float64
	TotalCO2Weight float64
	CostWeight     float64
	DecayRate      float64
}

// Config holds the configuration values for the Green Scheduling plugin.
type Config struct {
	TimeSeriesConfig      TimeSeriesConfig
	SustainabilityWeights SustainabilityWeights
	SerialNumLabel        string
}
