# Grafana Loki Stack

This directory contains the configuration files for running Grafana and Loki locally for testing the ionlog Loki integration.

## Quick Start

### Start the Stack
```bash
make loki-up
```

### Stop and Clean Up
```bash
make loki-down
```

### View Logs
```bash
make loki-logs
```

### Check Status
```bash
make loki-status
```

## Services

### Grafana
- **URL**: http://localhost:3000
- **Username**: admin
- **Password**: admin
- **Features**:
  - Pre-configured Loki data source
  - Ready-to-use dashboards
  - Log exploration and visualization

### Loki
- **URL**: http://localhost:3100
- **Features**:
  - Log aggregation and storage
  - Label-based indexing
  - REST API for log ingestion

## Configuration Files

### `loki-config.yml`
Loki server configuration with:
- File system storage
- Optimized for development
- Retention policies
- Performance tuning

### `grafana-datasources.yml`
Grafana data source configuration:
- Automatic Loki data source setup
- Default data source configuration
- Derived fields for trace correlation

## Usage with ionlog

Once the stack is running, you can use the ionlog Loki integration:

```go
package main

import (
    "github.com/IonicHealthUsa/ionlog"
    "github.com/IonicHealthUsa/ionlog/internal/observability/loki"
)

func main() {
    // Configure logger with Loki integration
    ionlog.SetAttributes(
        ionlog.WithLokiIntegration(
            loki.WithLoki(
                "http://localhost:3100",
                map[string]string{
                    "service":     "my-app",
                    "environment": "development",
                },
            ),
        ),
    )

    // Start logging
    ionlog.Start()
    defer ionlog.Stop()

    // Your logs will automatically be sent to Loki
    ionlog.Info("Hello from ionlog!")
}
```

## Querying Logs in Grafana

### Basic Queries
```
{service="my-app"}
{service="my-app", level="ERROR"}
{service="my-app"} |= "error"
```

### JSON Parsing
```
{service="my-app"} | json | level="INFO"
{service="my-app"} | json | duration > 100
```

### Time-based Queries
```
{service="my-app"} |= "error" | line_format "{{.timestamp}} {{.msg}}"
```

## Troubleshooting

### Check Container Status
```bash
make loki-status
```

### View Container Logs
```bash
make loki-logs
```

### Restart Services
```bash
make loki-down
make loki-up
```

### Access Loki API Directly
```bash
curl http://localhost:3100/ready
curl http://localhost:3100/metrics
```

## Data Persistence

- **Loki Data**: Stored in Docker volume `ionlog_loki-data`
- **Grafana Data**: Stored in Docker volume `ionlog_grafana-data`
- **Cleanup**: Use `make loki-down` to remove all data

## Performance Notes

- Configuration optimized for development
- File system storage (not production-ready)
- Single instance setup
- No clustering or high availability

For production deployments, refer to the official Grafana Loki documentation.
