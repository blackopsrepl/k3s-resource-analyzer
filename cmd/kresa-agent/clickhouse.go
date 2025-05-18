package main

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type ClickHouseClient struct {
	conn clickhouse.Conn
}

func newClickHouseClient(addr string) (*ClickHouseClient, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		return nil, err
	}
	return &ClickHouseClient{conn: conn}, nil
}

func (c *ClickHouseClient) close() error {
	return c.conn.Close()
}

func (c *ClickHouseClient) storeMetrics(ctx context.Context, metrics *ResourceMetrics) error {
	batch, err := c.conn.PrepareBatch(ctx, "INSERT INTO metrics (timestamp, type, name, namespace, cpu_usage, memory_usage)")
	if err != nil {
		return err
	}

	for _, pm := range metrics.PodMetrics {
		batch.Append(time.Now(), "pod", pm.PodName, pm.Namespace, pm.CPUUsage, pm.MemoryUsage)
	}

	for _, nm := range metrics.NodeMetrics {
		batch.Append(time.Now(), "node", nm.NodeName, "", nm.CPUUsage, nm.MemoryUsage)
	}

	return batch.Send()
}
