package greenscheduling

import (
	"context"
	"errors"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/scheduler-plugins/apis/config"
	"sigs.k8s.io/scheduler-plugins/pkg/greenscheduling/kubeinfo"
	"sigs.k8s.io/scheduler-plugins/pkg/greenscheduling/sicclient"
	"sigs.k8s.io/scheduler-plugins/pkg/greenscheduling/sicclient/sicparams"
	"sigs.k8s.io/scheduler-plugins/pkg/greenscheduling/sicclient/sicresponse"
	"sigs.k8s.io/scheduler-plugins/pkg/greenscheduling/sustainabilityprofile"
)

const (
	// Name is the name of the plugin used in Registry and configurations.
	Name = "GreenScheduling"

	// Scaling factor to prevent low scores from being rounded to 0.
	scoreScalingFactor = 1000
)

// GreenScheduling encapsulates the dependencies and configuration needed to execute
// sustainability-aware scheduling in Kubernetes clusters. By accessing external data
// on carbon emissions, energy usage, and costs, this struct enables the plugin to
// assess the environmental impact of node selection, supporting Kubernetes in making
// eco-conscious scheduling decisions.
type GreenScheduling struct {
	// config holds the configuration parameters needed for the plugin, such as
	// time series data settings, weighting factors for sustainability metrics,
	// and other options required to calculate sustainability scores effectively.
	config Config

	// kubeClient is a Kubernetes API client used to interact with Kubernetes resources,
	// retrieve node details, and obtain node labels, which can be used as identifiers
	// to map nodes with external sustainability data.
	kubeClient *kubeinfo.KubeClient

	// sicClient is a client for the Sustainability Information Center (SIC) API, which
	// provides environmental data such as CO2 emissions, energy consumption, and cost.
	// This data is used to calculate sustainability scores for each node in the cluster.
	sicClient *sicclient.Client
}

var _ = framework.ScorePlugin(&GreenScheduling{})

// New initializes a GreenScheduling plugin with the provided configuration arguments.
// It validates the input configuration, creates necessary clients, and constructs
// the `GreenScheduling` instance with all required dependencies and settings.
func New(_ context.Context, obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	// Attempt to cast the incoming object to GreenSchedulingArgs to access user-defined settings.
	args, ok := obj.(*config.GreenSchedulingArgs)
	if !ok {
		return nil, fmt.Errorf("expected GreenSchedulingArgs but got %T", obj)
	}

	// Validate provided arguments to ensure necessary configuration values are set correctly.
	if err := ValidateGreenSchedulingArgs(args); err != nil {
		return nil, fmt.Errorf("validation failed for GreenSchedulingArgs: %w", err)
	}

	// Initialize the Kubernetes client to allow interaction with Kubernetes resources.
	kubeClient, err := kubeinfo.NewKubeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client for GreenSchedulingArgs: %w", err)
	}

	// Define the token configuration to authenticate with the Sustainability Information Center (SIC).
	tokenConfig := sicclient.TokenConfig{
		URL:          args.TokenURL,
		ClientID:     args.ClientID,
		ClientSecret: args.ClientSecret,
	}

	// Initialize the SIC client with the provided hostname and token configuration
	// to retrieve environmental metrics like CO2 emissions and energy usage.
	sicClient := sicclient.New(sicclient.Config{
		Hostname:    args.SICHostname,
		TokenConfig: tokenConfig,
	})

	// Build the Config object that will encapsulate all plugin-specific settings,
	// including time series settings and sustainability metric weights.
	config := Config{
		TimeSeriesConfig: TimeSeriesConfig{
			SeriesInterval: args.TimeSeriesInterval, // Interval for time series data (e.g., hourly, daily)
			DaysToConsider: args.ConsiderationDays,  // Number of days to look back for sustainability data
		},
		SustainabilityWeights: SustainabilityWeights{
			args.CO2DecayWeight, // Weight for decaying CO2 emissions
			args.TotalCO2Weight, // Weight for total CO2 emissions
			args.CostWeight,     // Weight for cost in sustainability score calculation
			args.DecayRate,      // Rate at which CO2 impact decays over time
		},
		SerialNumLabel: args.SerialNumLabel, // Node label to identify the serial number for SIC lookup
	}

	// Return a new instance of GreenScheduling with all necessary clients and configurations.
	return &GreenScheduling{
		kubeClient: kubeClient, // Client for interacting with Kubernetes resources
		sicClient:  sicClient,  // Client for accessing environmental data from SIC
		config:     config,     // Plugin configuration settings
	}, nil
}

// Name returns name of the plugin. It is used in logs, etc.
func (n *GreenScheduling) Name() string {
	return Name
}

// Score computes the sustainability score for a given node.
func (gks *GreenScheduling) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	// Retrieve the node's serial number
	serialNum, err := gks.getNodeSerialNum(nodeName)
	if err != nil {
		if errors.Is(err, kubeinfo.ErrNodeNotFound) || errors.Is(err, kubeinfo.ErrLabelNotFound) {
			klog.Infof("Node %s missing required label '%s': %v", nodeName, gks.config.SerialNumLabel, err)
			return 0, framework.NewStatus(framework.Success)
		}
		klog.Errorf("Error retrieving label '%s' for node %s: %v", nodeName, gks.config.SerialNumLabel, err)
		return 0, framework.AsStatus(err)
	}

	// Build SIC parameters
	params, err := gks.buildSicParams(serialNum)
	if err != nil {
		klog.Errorf("Error creating filter parameters for node %s: %v", nodeName, err)
		return 0, framework.AsStatus(err)
	}

	// Fetch usage data from SIC
	usageByEntity, usageSeries, err := gks.fetchSicData(params)
	if err != nil {
		klog.Errorf("Error fetching SIC data for node %s: %v", nodeName, err)
		return 0, framework.AsStatus(err)
	}

	// Build emission data points
	dataPoints, err := gks.buildEmissionDataPoints(usageSeries)
	if err != nil {
		klog.Errorf("Error building emission data points for node %s: %v", nodeName, err)
		return 0, framework.AsStatus(err)
	}

	// Calculate the sustainability score
	score := gks.calculateSustainabilityScore(usageByEntity, dataPoints)
	klog.Infof("Calculated sustainability score for node %s with serial number %s: %f", nodeName, serialNum, score)

	// Scale the score to preserve precision and return as an integer
	scaledScore := int64(score * scoreScalingFactor)

	// Return the calculated score with a success status
	return scaledScore, framework.NewStatus(framework.Success)
}

// getNodeSerialNum retrieves the serial number label from the node.
func (gks *GreenScheduling) getNodeSerialNum(nodeName string) (string, error) {
	serialNum, err := gks.kubeClient.GetNodeLabelValue(nodeName, gks.config.SerialNumLabel)
	if err != nil {
		return "", err
	}

	return serialNum, nil
}

// buildSicParams constructs SIC parameters based on the serial number.
func (gks *GreenScheduling) buildSicParams(serialNum string) (*sicparams.Params, error) {
	filter, err := sicparams.NewFilter(sicparams.FilterKeyEntitySerialNum, sicparams.FilterOperatorEquals, serialNum)
	if err != nil {
		return sicparams.New(), fmt.Errorf("error creating filter: %w", err)
	}

	return sicparams.New().AddFilter(filter), nil
}

// fetchSicData retrieves usage data and time series data from the SIC client.
func (gks *GreenScheduling) fetchSicData(params *sicparams.Params) (*sicresponse.UsageByEntityResponse, *sicresponse.UsageSeriesResponse, error) {
	startTime := time.Now().AddDate(0, 0, -int(gks.config.TimeSeriesConfig.DaysToConsider)).Format(time.RFC3339)
	endTime := time.Now().Format(time.RFC3339)

	usageByEntity, err := gks.sicClient.GetUsageByEntity(startTime, endTime, params)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching usage by entity data: %w", err)
	}

	usageSeries, err := gks.sicClient.GetUsageSeries(startTime, endTime, gks.config.TimeSeriesConfig.SeriesInterval, params)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching usage series data: %w", err)
	}

	return usageByEntity, usageSeries, nil
}

// buildEmissionDataPoints processes usage series data into emission data points.
func (gks *GreenScheduling) buildEmissionDataPoints(usageSeries *sicresponse.UsageSeriesResponse) ([]sustainabilityprofile.EmissionDataPoint, error) {
	var dataPoints []sustainabilityprofile.EmissionDataPoint
	for _, entity := range usageSeries.Items {
		timeBucket, err := entity.GetTimeBucket()
		if err != nil {
			return nil, fmt.Errorf("error retrieving time bucket: %w", err)
		}
		dataPoints = append(dataPoints, sustainabilityprofile.EmissionDataPoint{
			Co2:  entity.GetCo2eMetricTon(),
			Time: timeBucket,
		})
	}
	return dataPoints, nil
}

// calculateSustainabilityScore calculates the sustainability score from entity data and emission data points.
func (gks *GreenScheduling) calculateSustainabilityScore(usageByEntity *sicresponse.UsageByEntityResponse, dataPoints []sustainabilityprofile.EmissionDataPoint) float64 {
	// Check if the usageByEntity response contains items
	if len(usageByEntity.Items) == 0 {
		klog.Warning("UsageByEntity response is empty")
		return 0.0
	}

	sData := sustainabilityprofile.New(
		dataPoints,
		usageByEntity.Items[0].GetCo2eMetricTon(),
		usageByEntity.Items[0].GetCostUsd(),
	)
	return sData.CalculateScore(
		sustainabilityprofile.NewSustainabilityWeights(
			gks.config.SustainabilityWeights.CO2DecayWeight,
			gks.config.SustainabilityWeights.TotalCO2Weight,
			gks.config.SustainabilityWeights.CostWeight,
			gks.config.SustainabilityWeights.DecayRate,
		),
	)
}

// ScoreExtensions of the GreenScheduling plugin to implement the framework.ScoreExtensions interface.
func (gks *GreenScheduling) ScoreExtensions() framework.ScoreExtensions {
	return gks
}

// NormalizeScore to scale every score to the framework.MaxNodeScore.
func (gks *GreenScheduling) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	var higherScore int64
	for _, node := range scores {
		if higherScore < node.Score {
			higherScore = node.Score
		}
	}

	for i, node := range scores {
		scores[i].Score = node.Score * framework.MaxNodeScore / higherScore
	}

	return nil
}
