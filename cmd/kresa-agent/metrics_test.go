package main

import (
	"context"
	"flag"
	"log"
	"testing"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCollectMetrics(t *testing.T) {
	var config *rest.Config

	var err error

	kubeconfig := flag.String("kubeconfig", "/home/pvd/.kube/config", "Path to kubeconfig file (if not in-cluster)")

	config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	collector, err := newMetricsCollector(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	for i := 0; i <= 10; i++ {
		// Collect metrics
		metrics, _ := collector.collectMetrics(ctx)
		log.Print(metrics)
	}

}
