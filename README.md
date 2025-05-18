# k3s-resource-analyzer

- kresa-agent: in-cluster metrics collector
- kresa: cli-tool

To create the required Clickhouse table:
```sql
CREATE TABLE default.metrics (timestamp DateTime, type String, name String, namespace String, cpu_usage Int64, memory_usage Int64) ENGINE = MergeTree() ORDER BY (timestamp, type, name)
```
