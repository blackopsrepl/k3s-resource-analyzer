package main

import (
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
