# Loki Integration for ionlog

This document describes how to integrate Grafana Loki with the ionlog library for centralized log aggregation and analysis.

## Public API

The Loki integration provides a **completely public API** - no need to import internal packages! All configuration and management functions are available directly from the main `ionlog` package.

### Quick Start

```go
package main

import (
    "github.com/IonicHealthUsa/ionlog"
)

func main() {
    // Configure Loki integration using public API
    config := ionlog.WithLoki("http://localhost:3100", map[string]string{
        "service": "my-app",
        "version": "1.0.0",
    })
    
    // Set up the logger
    ionlog.SetAttributes(ionlog.WithLokiIntegration(config))
    ionlog.Start()
    defer ionlog.Stop() // Graceful shutdown is automatic
    
    // Your application code
    ionlog.Info("Application started")
}
```

### Key Benefits

- ✅ **No internal imports** - Everything is in the main `ionlog` package
- ✅ **Clean API** - Builder pattern for configuration
- ✅ **Transparent shutdown** - No warnings or errors exposed
- ✅ **Configurable timeouts** - Environment variables and code-level settings
- ✅ **Type safety** - Strongly typed configuration structs

## Overview

The Loki integration provides:
- **Asynchronous log forwarding** to Grafana Loki
- **Label-based indexing** for fast log queries
- **Batch processing** for optimal performance
- **Authentication support** for secure deployments
- **Multi-tenant support** for enterprise environments
- **Configurable timeouts and batch sizes**

## Quick Start

### Basic Integration

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

    // Start the logger
    ionlog.Start()
    defer ionlog.Stop()

    // Log messages (automatically sent to Loki)
    ionlog.Info("Application started")
    ionlog.Error("Database connection failed")
}
```

### Advanced Configuration

```go
import (
    "time"
    "github.com/IonicHealthUsa/ionlog"
    "github.com/IonicHealthUsa/ionlog/internal/observability/loki"
)

func main() {
    // Create custom configuration
    config := loki.WithLoki("http://localhost:3100", map[string]string{
        "service":     "my-app",
        "environment": "production",
        "version":     "1.0.0",
    })

    // Add authentication
    config = loki.WithLokiAuth(config, "admin", "password")

    // Add tenant ID for multi-tenancy
    config = loki.WithLokiTenant(config, "tenant-123")

    // Set custom batch size and timeout
    config = loki.WithLokiBatchSize(config, 100)
    config = loki.WithLokiTimeout(config, 30*time.Second)

    // Configure logger with advanced Loki integration
    ionlog.SetAttributes(
        ionlog.WithLokiIntegration(config),
    )

    // Start logging
    ionlog.Start()
    defer ionlog.Stop()

    ionlog.Info("Application started with advanced Loki configuration")
}
```

## Graceful Shutdown

The Loki integration includes graceful shutdown functionality to ensure all buffered logs are sent to Loki before the application exits.

### Automatic Graceful Shutdown

When you call `ionlog.Stop()`, the library automatically attempts to gracefully shutdown the Loki integration with a configurable timeout (default: 10 seconds).

```go
// The Stop() function automatically handles graceful Loki shutdown
ionlog.Stop()
```

### Custom Shutdown Timeout

You can configure the shutdown timeout:

```go
ionlog.SetAttributes(
    ionlog.WithLokiIntegration(config),
    ionlog.WithLokiShutdownTimeout(5*time.Second), // 5 second timeout
)
```

### Manual Graceful Shutdown

For more control, you can manually shutdown the Loki integration:

```go
// Get the integration and shutdown manually
if integration := ionlog.GetLokiIntegration(); integration != nil {
    err := integration.GracefulShutdown(5 * time.Second)
    if err != nil {
        log.Printf("Loki shutdown failed: %v", err)
    }
}
```

### Signal Handling Example

```go
package main

import (
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/IonicHealthUsa/ionlog"
    "github.com/IonicHealthUsa/ionlog/internal/observability/loki"
)

func main() {
    // Configure with custom shutdown timeout
    ionlog.SetAttributes(
        ionlog.WithLokiIntegration(
            loki.WithLoki("http://localhost:3100", map[string]string{
                "service": "my-app",
            }),
        ),
        ionlog.WithLokiShutdownTimeout(5*time.Second),
    )
    
    ionlog.Start()
    defer ionlog.Stop() // This will gracefully shutdown Loki
    
    // Set up signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Your application logic here...
    
    // Wait for shutdown signal
    <-sigChan
    fmt.Println("Shutting down gracefully...")
    // ionlog.Stop() will be called by defer, ensuring logs are sent
}
```

### Shutdown Behavior

- **Success**: All buffered logs are sent to Loki within the timeout
- **Timeout**: If logs cannot be sent within the timeout, the application will exit with a warning
- **No Integration**: If no Loki integration is configured, shutdown completes immediately

## Configuration Options

### Environment Variables

The Loki integration supports configuration via environment variables:

```bash
# Loki Configuration
export LOKI_URL=http://localhost:3100
export LOKI_USERNAME=admin
export LOKI_PASSWORD=admin123
export LOKI_TENANT_ID=tenant-123
export LOKI_BATCH_SIZE=100
export LOKI_TIMEOUT=30s

# Service Labels
export SERVICE_NAME=my-app
export ENVIRONMENT=production
export VERSION=1.0.0
export COMPONENT=logger
```

### Configuration Structure

```go
type LokiConfig struct {
    URL       string            // Loki server URL
    Username  string            // Basic auth username
    Password  string            // Basic auth password
    TenantID  string            // Multi-tenant ID
    Labels    map[string]string // Default labels
    BatchSize int               // Batch size for sending logs
    Timeout   time.Duration     // Request timeout
}
```

## Label Strategy

### Default Labels

The integration automatically adds these labels to all log entries:

- `service`: Service name (from SERVICE_NAME env var)
- `environment`: Environment (from ENVIRONMENT env var)
- `version`: Application version (from VERSION env var)
- `component`: Component name (from COMPONENT env var)

### Log-Specific Labels

Each log entry automatically includes:

- `level`: Log level (INFO, ERROR, WARN, DEBUG, TRACE)
- `file`: Source file name
- `package`: Go package name
- `function`: Function name
- `line`: Line number

### Custom Labels

You can add custom labels through configuration:

```go
config := loki.WithLoki("http://localhost:3100", map[string]string{
    "service":     "my-app",
    "environment": "production",
    "team":        "platform",
    "region":      "us-east-1",
    "datacenter":  "dc1",
    "cluster":     "prod-cluster",
})
```

## Query Examples

### Basic Queries

```logql
# All logs from a service
{service="my-app"}

# Error logs only
{service="my-app", level="ERROR"}

# Logs from specific environment
{service="my-app", environment="production"}

# Logs containing specific text
{service="my-app"} |= "error"
```

### Advanced Queries

```logql
# JSON parsing for structured logs
{service="my-app"} | json | level="INFO"

# Filter by duration
{service="my-app"} | json | duration > 100

# Filter by user ID
{service="my-app"} | json | user_id="12345"

# Rate of errors over time
rate({service="my-app", level="ERROR"}[5m])

# Top error messages
topk(10, count_over_time({service="my-app", level="ERROR"}[1h]))
```

### Performance Queries

```logql
# Average response time
avg_over_time({service="my-app"} | json | unwrap duration [5m])

# 95th percentile response time
quantile_over_time(0.95, {service="my-app"} | json | unwrap duration [5m])

# Request rate
rate({service="my-app"} | json | operation="request" [1m])
```

## Performance Considerations

### Batch Processing

The integration uses batch processing to optimize performance:

- **Default batch size**: 100 logs
- **Automatic flushing**: Every 5 seconds
- **Manual flushing**: When buffer is full
- **Configurable batch size**: Adjust based on your needs

### Memory Usage

- **Buffer size**: Configurable (default 100 entries)
- **Background flushing**: Prevents memory buildup
- **Graceful shutdown**: Flushes remaining logs

### Network Optimization

- **Compression**: Logs are compressed before sending
- **Connection pooling**: Reuses HTTP connections
- **Timeout handling**: Configurable request timeouts

## Error Handling

### Connection Failures

The integration handles connection failures gracefully:

- **Retry logic**: Built into the HTTP client
- **Non-blocking**: Log failures don't block application
- **Error logging**: Connection errors are logged to stderr

### Rate Limiting

If Loki is rate limiting requests:

- **Backoff**: Automatic backoff on rate limit errors
- **Batch reduction**: Reduces batch size on errors
- **Circuit breaker**: Stops sending on repeated failures

## Monitoring and Observability

### Health Checks

```go
integration := ionlog.GetLokiIntegration()
if integration != nil {
    err := integration.HealthCheck()
    if err != nil {
        log.Printf("Loki health check failed: %v", err)
    }
}
```

### Statistics

```go
integration := ionlog.GetLokiIntegration()
if integration != nil {
    stats := integration.GetStats()
    log.Printf("Loki stats: %+v", stats)
}
```

### Metrics

The integration provides these metrics:

- `loki_logs_sent_total`: Total logs sent to Loki
- `loki_logs_failed_total`: Total failed log sends
- `loki_batch_size`: Current batch size
- `loki_buffer_size`: Current buffer size

## Best Practices

### Label Design

1. **Keep cardinality low**: Avoid high-cardinality labels
2. **Use consistent naming**: Follow naming conventions
3. **Include essential context**: Add labels that help with filtering
4. **Avoid sensitive data**: Don't include passwords or tokens

### Performance

1. **Tune batch size**: Adjust based on log volume
2. **Monitor memory usage**: Watch buffer size
3. **Use appropriate timeouts**: Balance reliability vs performance
4. **Consider log levels**: Use appropriate log levels

### Security

1. **Use authentication**: Always use auth in production
2. **Secure transport**: Use HTTPS for Loki communication
3. **Rotate credentials**: Regularly rotate auth credentials
4. **Monitor access**: Log and monitor Loki access

## Troubleshooting

### Common Issues

#### Logs Not Appearing in Loki

1. Check Loki server is running
2. Verify URL and authentication
3. Check network connectivity
4. Review Loki server logs

#### High Memory Usage

1. Reduce batch size
2. Increase flush frequency
3. Check for memory leaks
4. Monitor buffer size

#### Slow Performance

1. Increase batch size
2. Reduce flush frequency
3. Check network latency
4. Monitor Loki server performance

### Debug Mode

Enable debug logging:

```go
// Set trace mode for detailed logging
ionlog.SetAttributes(ionlog.WithTraceMode(true))
```

### Log Analysis

Use these queries to debug issues:

```logql
# Check if logs are being received
{service="my-app"} | json | level="ERROR"

# Monitor log volume
rate({service="my-app"}[1m])

# Check for connection errors
{service="my-app"} |= "connection failed"
```

## Transparent Graceful Shutdown

The Loki integration includes a **completely transparent** graceful shutdown mechanism that ensures all buffered logs are sent to Loki before the application exits, with configurable timeout protection.

### Key Features

- **Completely Transparent**: No warnings or errors are exposed to the user
- **Configurable Timeout**: Set timeout via environment variable or code
- **Automatic Integration**: Works automatically with `ionlog.Stop()`
- **Timeout Protection**: Prevents application hanging
- **Priority-based Configuration**: Code > Environment > Default

### Configuration Priority

The shutdown timeout is configured in the following priority order:

1. **Programmatically set** (via `WithLokiShutdownTimeout`)
2. **Environment variable** (`IONLOG_LOKI_SHUTDOWN_TIMEOUT`)
3. **Default timeout** (10 seconds)

### Environment Variable Configuration

```bash
# Set timeout via environment variable
export IONLOG_LOKI_SHUTDOWN_TIMEOUT="5s"
export IONLOG_LOKI_SHUTDOWN_TIMEOUT="30s"
export IONLOG_LOKI_SHUTDOWN_TIMEOUT="1m"
```

### Code Configuration

```go
// Set custom timeout in code (overrides environment)
ionlog.SetAttributes(
    ionlog.WithLokiIntegration(config),
    ionlog.WithLokiShutdownTimeout(5*time.Second),
)

// Reset to use environment variable or default
ionlog.ResetLokiShutdownTimeout()

// Get current timeout
timeout := ionlog.GetLokiShutdownTimeout()
```

### Automatic Integration

The graceful shutdown is **completely automatic** and integrated with `ionlog.Stop()`:

```go
// This automatically handles Loki graceful shutdown - no manual intervention needed
ionlog.Stop() // Completely transparent, no warnings, no errors
```

**No manual shutdown required** - the graceful shutdown happens automatically when you call `ionlog.Stop()`.

### Timeout Protection

- **Minimum timeout**: 1 second (prevents immediate failures)
- **Default timeout**: 10 seconds
- **Graceful timeout**: If timeout is reached, shutdown continues without errors
- **No application hanging**: Application always exits within the specified timeout

### Example Usage

```go
package main

import (
    "time"
    "github.com/IonicHealthUsa/ionlog"
    "github.com/IonicHealthUsa/ionlog/internal/observability/loki"
)

func main() {
    // Setup with custom timeout
    config := loki.WithLoki("http://localhost:3100")
    ionlog.SetAttributes(
        ionlog.WithLokiIntegration(config),
        ionlog.WithLokiShutdownTimeout(5*time.Second),
    )
    
    // Your application code
    ionlog.Info("Application running")
    
    // Graceful shutdown - completely transparent
    ionlog.Stop() // All logs sent to Loki, no warnings
}
```

## Examples

See the following examples for complete usage:

- [Basic Example](examples/loki/basic/main.go)
- [Advanced Example](examples/loki/advanced/main.go)
- [Graceful Shutdown Example](examples/loki/graceful_shutdown/main.go)
- [Transparent Shutdown Example](examples/loki/transparent_shutdown/main.go)

## API Reference

### Functions

- `WithLokiIntegration(config LokiConfig) customAttrs` - Configure Loki integration
- `GetLokiIntegration() *LokiIntegration` - Get current Loki integration
- `WithLokiShutdownTimeout(timeout time.Duration) customAttrs` - Set shutdown timeout
- `GetLokiShutdownTimeout() time.Duration` - Get current shutdown timeout
- `ResetLokiShutdownTimeout()` - Reset timeout to environment/default

### Configuration Builders

- `WithLoki(url string, labels map[string]string) LokiConfig` - Create basic Loki configuration
- `WithLokiAuth(config LokiConfig, username, password string) LokiConfig` - Add authentication
- `WithLokiTenant(config LokiConfig, tenantID string) LokiConfig` - Add tenant ID
- `WithLokiBatchSize(config LokiConfig, batchSize int) LokiConfig` - Set batch size
- `WithLokiTimeout(config LokiConfig, timeout time.Duration) LokiConfig` - Set timeout

## Support

For issues and questions:

1. Check the troubleshooting section
2. Review the examples
3. Check Loki server documentation
4. Open an issue in the repository
