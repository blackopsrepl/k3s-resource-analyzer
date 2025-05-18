package main

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// Collects resource utilization metrics from a K3s cluster
type MetricsCollector struct {
	metricsClient *versioned.Clientset
}

// Represents pod and node metrics
type ResourceMetrics struct {
	PodMetrics  []PodMetric
	NodeMetrics []NodeMetric
}

// Represents metrics for single pod
type PodMetric struct {
	Namespace   string
	PodName     string
	CPUUsage    int64 // in millicores
	MemoryUsage int64 // in megabytes
}

// Represents metrics for a single node
type NodeMetric struct {
	NodeName       string
	CPUUsage       int64 // in millicores
	MemoryUsage    int64 // in megabytes
	CPUCapacity    int64 // in millicores
	MemoryCapacity int64 // in megabytes
}

func newMetricsCollector(config *rest.Config) (*MetricsCollector, error) {
	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create metrics client: %v", err)
	}

	return &MetricsCollector{metricsClient: metricsClient}, nil
}

func (collector *MetricsCollector) collectMetrics(ctx context.Context) (*ResourceMetrics, error) {
	metrics := &ResourceMetrics{}
	var err error

	// Use channels for concurrent collection
	podChan := make(chan []PodMetric)
	nodeChan := make(chan []NodeMetric)
	errChan := make(chan error, 2)

	// Run goroutines
	go collector.collectPodMetrics(ctx, podChan, errChan)
	go collector.collectNodeMetrics(ctx, nodeChan, errChan)

	select {
	case metrics.PodMetrics = <-podChan:
	case err = <-errChan:
		return nil, err
	}
	select {
	case metrics.NodeMetrics = <-nodeChan:
	case err = <-errChan:
		return nil, err
	}

	return metrics, nil
}

func (collector *MetricsCollector) collectPodMetrics(ctx context.Context, podChan chan<- []PodMetric, errChan chan<- error) {
	podMetricsList, err := collector.metricsClient.MetricsV1beta1().PodMetricses("").List(ctx, v1.ListOptions{})
	if err != nil {
		errChan <- fmt.Errorf("Failed to list pod metrics: %v", err)
		return
	}
	var podMetrics []PodMetric
	for _, pm := range podMetricsList.Items {
		for _, container := range pm.Containers {
			podMetrics = append(podMetrics, PodMetric{
				Namespace:   pm.Namespace,
				PodName:     pm.Name,
				CPUUsage:    container.Usage.Cpu().MilliValue(),
				MemoryUsage: container.Usage.Memory().Value() / (1024 * 1024),
			})
		}
	}
	podChan <- podMetrics
}

func (collector *MetricsCollector) collectNodeMetrics(ctx context.Context, nodeChan chan<- []NodeMetric, errChan chan<- error) {
	nodeMetricsList, err := collector.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, v1.ListOptions{})
	if err != nil {
		errChan <- fmt.Errorf("Failed to list node metrics: %v", err)
		return
	}
	var nodeMetrics []NodeMetric
	for _, nm := range nodeMetricsList.Items {
		nodeMetrics = append(nodeMetrics, NodeMetric{
			NodeName:       nm.Name,
			CPUUsage:       nm.Usage.Cpu().MilliValue(),
			MemoryUsage:    nm.Usage.Memory().Value() / (1024 * 1024),
			CPUCapacity:    6000,
			MemoryCapacity: 65536,
		})
	}
	nodeChan <- nodeMetrics
}
