package main

import (
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type MetricsCollector struct {
	metricsClient *versioned.Clientset
}
1
type ResourceMetrics struct {
	PodMetrics []PodMetric
	NodesMetrics []NodeMetric
}

type PodMetric struct {}

type NodeMetric struct {}
