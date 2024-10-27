package kubeinfo

import (
	"context"
	"errors"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Predefined errors for specific cases
var (
	ErrNodeNotFound   = errors.New("node not found")
	ErrLabelNotFound  = errors.New("label not found on node")
	ErrClientCreation = errors.New("failed to create Kubernetes client")
)

// KubeClient wraps a Kubernetes client to interact with the cluster.
type KubeClient struct {
	clientset *kubernetes.Clientset
}

// NewKubeClient creates a new KubeClient using in-cluster configuration.
func NewKubeClient() (*KubeClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrClientCreation, err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrClientCreation, err)
	}

	return &KubeClient{clientset: clientset}, nil
}

// GetNodeLabels retrieves the labels of a specified node by node name.
func (kc *KubeClient) GetNodeLabels(nodeName string) (map[string]string, error) {
	node, err := kc.clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, ErrNodeNotFound
		}

		return nil, fmt.Errorf("failed to get node %s: %w", nodeName, err)
	}
	return node.Labels, nil
}

// GetNodeLabelValue retrieves the value of a specific label for a given node.
func (kc *KubeClient) GetNodeLabelValue(nodeName, label string) (string, error) {
	node, err := kc.clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return "", ErrNodeNotFound
		}

		return "", fmt.Errorf("failed to get node %s: %w", nodeName, err)
	}

	value, exists := node.Labels[label]
	if !exists {
		return "", ErrLabelNotFound
	}

	return value, nil
}
