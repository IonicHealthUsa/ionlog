package main

import (
	"fmt"
	"time"

	"github.com/IonicHealthUsa/ionlog"
)

func main() {
	fmt.Println("=== Loki Integration Example ===")

	// Configure logger with Loki integration
	ionlog.SetAttributes(
		ionlog.WithLokiIntegration(
			ionlog.WithLoki(
				"http://localhost:3100",
				map[string]string{
					"service":     "loki-example",
					"environment": "development",
					"component":   "logger",
				},
			),
		),
	)

	// Start the logger
	ionlog.Start()
	defer ionlog.Stop()

	// Add static fields to the logger
	ionlog.SetAttributes(
		ionlog.WithStaticFields(map[string]string{
			"service":     "loki-example",
			"environment": "development",
			"version":     "1.0.0",
		}),
	)

	fmt.Println("\n1. Basic Logging to Loki")
	ionlog.Info("This is a basic info log message for Loki.")
	ionlog.Debug("This is a debug message with some context.")

	fmt.Println("\n2. Error Logging to Loki")
	ionlog.Error("An unexpected error occurred during processing.")
	ionlog.Warn("A warning message indicating a potential issue.")

	fmt.Println("\n3. Different Log Levels")
	ionlog.Trace("Entering function 'processData'")
	ionlog.Info("User 'john.doe' logged in successfully.")
	ionlog.Warn("Disk space is running low on server 'app-server-01'.")
	ionlog.Error("Failed to connect to external API: connection refused.")

	fmt.Println("\n4. LogOnce Functionality")
	ionlog.LogOnceInfo("This info message should only appear once.")
	ionlog.LogOnceInfo("This info message should only appear once.") // This will be ignored
	ionlog.LogOnceError("Critical error: database connection lost!")
	ionlog.LogOnceError("Critical error: database connection lost!") // This will be ignored

	fmt.Println("\n5. Business Logic Logging")
	ionlog.Info("Order processed successfully")
	ionlog.Info("Payment completed")
	ionlog.Warn("Low inventory warning")
	ionlog.Error("Payment gateway timeout")

	fmt.Println("\n6. Performance Logging")
	start := time.Now()
	time.Sleep(100 * time.Millisecond) // Simulate work
	ionlog.Info(fmt.Sprintf("Database query completed - duration: %dms, rows: %d", time.Since(start).Milliseconds(), 1500))

	start = time.Now()
	time.Sleep(50 * time.Millisecond) // Simulate work
	ionlog.Info(fmt.Sprintf("API request processed - duration: %dms, status: %d", time.Since(start).Milliseconds(), 200))

	fmt.Println("\n7. Security Logging")
	ionlog.Info("User login attempt - user_id: 12345, ip: 192.168.1.100")
	ionlog.Warn("Failed login attempt - user_id: 12345, ip: 192.168.1.100, reason: invalid_password")
	ionlog.Info("User logout - user_id: 12345")

	fmt.Println("\n8. Application Lifecycle Logging")
	ionlog.Info("Application starting up")
	ionlog.Info("Configuration loaded")
	ionlog.Info("Database connected")
	ionlog.Info("Server listening on port 8080")
	ionlog.Info("Application ready to accept requests")

	fmt.Println("\n9. High-Volume Logging Simulation")
	for range 10 {
		ionlog.Info(fmt.Sprintf("Processing request - request_id: req-%d, user_id: user-%d", time.Now().UnixNano()%1000, time.Now().UnixNano()%5))
	}

	fmt.Println("\n10. Metrics and Monitoring")
	ionlog.Info("System metrics - cpu_usage: 45.2%, memory_usage: 78.5%, disk_usage: 23.1%")
	ionlog.Info("Application metrics - requests_per_second: 150, response_time: 250ms, error_rate: 0.02")

	fmt.Println("\n=== Loki Example Completed ===")
	fmt.Println("✅ All logs automatically sent to Loki ASYNCHRONOUSLY!")
	fmt.Println()
	fmt.Println("🔍 Check Loki at: http://localhost:3100")
	fmt.Println("📊 Grafana: http://localhost:3000 (admin/admin)")
	fmt.Println()
	fmt.Println("📊 Query examples in Grafana:")
	fmt.Println(`  - {service="loki-example"}`)
	fmt.Println(`  - {service="loki-example", level="ERROR"}`)
	fmt.Println(`  - {service="loki-example"} |= "error"`)
	fmt.Println(`  - {service="loki-example"} | json | level="INFO"`)
	fmt.Println()
	fmt.Println("✨ Features demonstrated:")
	fmt.Println("  ✓ Automatic async log forwarding to Loki")
	fmt.Println("  ✓ All log levels (Debug, Info, Warn, Error, Trace)")
	fmt.Println("  ✓ LogOnce functionality")
	fmt.Println("  ✓ Static fields included in all logs")
	fmt.Println("  ✓ Caller information preserved")
	fmt.Println("  ✓ High performance - non-blocking")
	fmt.Println("  ✓ Label-based indexing for fast queries")
	fmt.Println("  ✓ Business logic logging")
	fmt.Println("  ✓ Performance monitoring")
	fmt.Println("  ✓ Security event logging")
	fmt.Println("  ✓ Application lifecycle logging")
	fmt.Println("  ✓ High-volume logging simulation")
	fmt.Println("  ✓ Metrics and monitoring")
	fmt.Println()
}
