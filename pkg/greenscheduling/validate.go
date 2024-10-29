package greenscheduling

import (
	"errors"

	"sigs.k8s.io/scheduler-plugins/apis/config"
)

// ValidateGreenSchedulingArgs checks if all required args are present.
func ValidateGreenSchedulingArgs(args *config.GreenSchedulingArgs) error {
	if args.TokenURL == "" || args.ClientID == "" || args.ClientSecret == "" || args.SICHostname == "" || args.SerialNumLabel == "" {
		return errors.New("missing required arguments")
	}
	if args.CO2DecayWeight < 0 || args.TotalCO2Weight < 0 || args.CostWeight < 0 || args.DecayRate < 0 || args.DecayRate > 1 {
		return errors.New("invalid weight or decay rate")
	}
	if args.TimeSeriesInterval == "" || args.ConsiderationDays <= 0 {
		return errors.New("invalid interval or consideration days")
	}
	return nil
}
