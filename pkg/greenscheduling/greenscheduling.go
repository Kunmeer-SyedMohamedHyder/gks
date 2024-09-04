package greenscheduling

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/scheduler-plugins/apis/config"
)

type GreenScheduling struct{}

var _ = framework.ScorePlugin(&GreenScheduling{})

func New(_ context.Context, obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	args, ok := obj.(*config.GreenSchedulingArgs)
	if !ok {
		return nil, fmt.Errorf("[GreenScheduling] want args to be of type GreenSchedulingArgs, got %T", obj)
	}

	klog.Infof("[GreenScheduling] args received. Sample: %s", args.Sample)
	return &GreenScheduling{}, nil
}

const (
	// Name is the name of the plugin used in Registry and configurations.
	Name = "GreenScheduling"
)

// Name returns name of the plugin. It is used in logs, etc.
func (n *GreenScheduling) Name() string {
	return Name
}

func (gks *GreenScheduling) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	var score int64
	if nodeName == "k8s-worker1" {
		score = 100
	}

	klog.Infof("[GreenScheduling] node '%s' score: %d", nodeName, score)
	return score, nil
}

func (gks *GreenScheduling) ScoreExtensions() framework.ScoreExtensions {
	return gks
}

func (gks *GreenScheduling) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	var higherScore int64
	for _, node := range scores {
		if higherScore < node.Score {
			higherScore = node.Score
		}
	}

	for i, node := range scores {
		scores[i].Score = framework.MaxNodeScore - (node.Score * framework.MaxNodeScore / higherScore)
	}

	klog.Infof("[GreenScheduling] Nodes final score: %v", scores)
	return nil
}
