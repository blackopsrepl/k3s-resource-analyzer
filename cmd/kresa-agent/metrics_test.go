package main

import (
	"context"
	"log"
	"os"
	"testing"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCollectMetrics(t *testing.T) {
	var config *rest.Config

	var err error

	envFile := "../../.env"

	setEnv(envFile)

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		log.Fatal(err)
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("KUBECONFIG must be set as environment variable!")
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
