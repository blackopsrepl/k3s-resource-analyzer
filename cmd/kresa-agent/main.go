package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// ENVIRONMENT //

	// Parse flags and load environment
	envFile := flag.String("env", "", "Path to .env file")
	flag.Parse()

	setEnv(*envFile)

	clickhouseAddr := os.Getenv("CLICKHOUSE_ADDR")
	log.Println(clickhouseAddr)
	if clickhouseAddr == "" {
		log.Fatalf("CLICKHOUSE_ADDR must be set as environment variable!")
	}
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		log.Fatalf("KUBECONFIG must be set as environment variable!")
	}

	flag.Parse()

	setEnv(*envFile)

	// CLIENTS //

	// Initialize Kubernetes client
	var config *rest.Config
	var err error

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to build config: %v", err)
	}

	// Initialize ClickHouse client
	chClient, err := newClickHouseClient(clickhouseAddr)
	if err != nil {
		log.Fatalf("Failed to initialize ClickHouse client: %v", err)
	}

	// Initialize metrics collector
	collector, err := newMetricsCollector(config)
	if err != nil {
		log.Fatalf("Failed to initialize collector: %v", err)
	}

	// Main loop
	ctx := context.Background()
	for {
		// Collect metrics
		metrics, err := collector.collectMetrics(ctx)
		if err != nil {
			log.Printf("Failed to collect metrics: %v", err)
			time.Sleep(10 * time.Second)
		}

		// Store metrics
		if err := chClient.storeMetrics(ctx, metrics); err != nil {
			log.Printf("Failed to store metrics in ClickHouse: %v", err)
		}
	}
}

func setEnv(envFile string) {
	if envFile != "" {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	} else if os.Getenv("CLICKHOUSE_ADDR") == "" || os.Getenv("KUBECONFIG") == "" {
		log.Fatalf("CLICKHOUSE_ADDR, KUBECONFIG must be set as environment variables!\n\nLoaded CLICKHOUSE_ADDR: %s, KUBECONFIG: %s", os.Getenv("CLICKHOUSE_ADDR"), os.Getenv("KUBECONFIG"))
	}
}
