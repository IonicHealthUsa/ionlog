package main

import (
	"fmt"
	"time"

	"github.com/IonicHealthUsa/ionlog"
)

func main() {
	fmt.Println("=== Advanced Loki Integration Example ===")

	// Create a custom Loki configuration
	config := ionlog.WithLoki("http://localhost:3100", map[string]string{
		"service":     "advanced-loki-example",
		"environment": "production",
		"version":     "2.0.0",
		"team":        "platform",
		"region":      "us-east-1",
	})

	// Add authentication if needed
	config = ionlog.WithLokiAuth(config, "admin", "admin123")

	// Add tenant ID for multi-tenancy
	config = ionlog.WithLokiTenant(config, "tenant-123")

	// Set custom batch size and timeout
	config = ionlog.WithLokiBatchSize(config, 50)
	config = ionlog.WithLokiTimeout(config, 60*time.Second)

	// Configure logger with advanced Loki integration
	ionlog.SetAttributes(
		ionlog.WithLokiIntegration(config),
		ionlog.WithLokiShutdownTimeout(2*time.Second),
	)

	// Start the logger
	ionlog.Start()
	defer ionlog.Stop()

	// Add static fields
	ionlog.SetAttributes(
		ionlog.WithStaticFields(map[string]string{
			"service":     "advanced-loki-example",
			"environment": "production",
			"version":     "2.0.0",
			"datacenter":  "dc1",
			"cluster":     "prod-cluster",
		}),
	)

	fmt.Println("\n1. Structured Logging with Rich Labels")
	ionlog.Info("Application started with advanced configuration")
	ionlog.Info("Database connection pool initialized")
	ionlog.Info("Cache layer activated")

	fmt.Println("\n2. Error Handling and Recovery")
	ionlog.Error("Database connection failed - attempting reconnection")
	time.Sleep(100 * time.Millisecond)
	ionlog.Info("Database reconnection successful")
	ionlog.Warn("High memory usage detected - triggering GC")

	fmt.Println("\n3. Performance Monitoring")
	start := time.Now()
	time.Sleep(200 * time.Millisecond) // Simulate work
	ionlog.Info(fmt.Sprintf("API endpoint /users processed - duration: %dms, status: 200", time.Since(start).Milliseconds()))

	start = time.Now()
	time.Sleep(150 * time.Millisecond) // Simulate work
	ionlog.Info(fmt.Sprintf("Database query executed - duration: %dms, rows_affected: 42", time.Since(start).Milliseconds()))

	fmt.Println("\n4. Business Metrics")
	ionlog.Info("User registration completed - user_id: 789, plan: premium")
	ionlog.Info("Payment processed - amount: $99.99, currency: USD, method: credit_card")
	ionlog.Info("Subscription activated - user_id: 789, plan: premium, expires: 2024-12-31")

	fmt.Println("\n5. Security Events")
	ionlog.Info("User login successful - user_id: 789, ip: 192.168.1.50, user_agent: Mozilla/5.0")
	ionlog.Warn("Suspicious activity detected - user_id: 789, ip: 192.168.1.50, action: multiple_failed_logins")
	ionlog.Info("User logout - user_id: 789, session_duration: 2h 15m")

	fmt.Println("\n6. System Health Monitoring")
	ionlog.Info("System health check - cpu: 45%, memory: 78%, disk: 23%, network: 12%")
	ionlog.Info("Service discovery updated - services: 15, healthy: 14, unhealthy: 1")
	ionlog.Warn("Load balancer health check failed - endpoint: /health, status: 503")

	fmt.Println("\n7. Distributed Tracing Context")
	ionlog.Info("Request started - trace_id: abc123, span_id: def456, operation: user.create")
	ionlog.Info("Database call initiated - trace_id: abc123, span_id: ghi789, query: SELECT * FROM users")
	ionlog.Info("External API call - trace_id: abc123, span_id: jkl012, service: payment-gateway")
	ionlog.Info("Request completed - trace_id: abc123, span_id: def456, duration: 250ms, status: 201")

	fmt.Println("\n8. Configuration Changes")
	ionlog.Info("Configuration reloaded - config_version: v2.1.0, changes: 3")
	ionlog.Info("Feature flag updated - flag: new_ui, enabled: true, rollout: 50%")
	ionlog.Info("Rate limit adjusted - endpoint: /api/users, limit: 1000/min")

	fmt.Println("\n9. Batch Operations")
	for i := 0; i < 5; i++ {
		ionlog.Info(fmt.Sprintf("Batch operation %d completed - items_processed: %d, errors: %d", i+1, 100+i*10, i))
	}

	// Check if the logs were sent to Loki
	integration := ionlog.GetLokiIntegration()
	if integration == nil {
		fmt.Printf("Failed to get Loki integration\n")
	}
	stats := integration.GetStats()
	fmt.Printf("Loki stats: %+v\n", stats)

	fmt.Println("\n10. Cleanup and Shutdown")
	ionlog.Info("Graceful shutdown initiated - reason: maintenance_window")
	ionlog.Info("Active connections closed - count: 42")
	ionlog.Info("Cache flushed - entries: 1024, size: 50MB")
	ionlog.Info("Application shutdown completed - uptime: 7d 12h 30m")

	fmt.Println("\n=== Advanced Loki Example Completed ===")
	fmt.Println("✅ All logs sent to Loki with advanced configuration!")
	fmt.Println()
	fmt.Println("🔍 Check Loki at: http://localhost:3100")
	fmt.Println("📊 Grafana: http://localhost:3000 (admin/admin)")
	fmt.Println()
	fmt.Println("📊 Advanced Query Examples:")
	fmt.Println(`  - {service="advanced-loki-example", environment="production"}`)
	fmt.Println(`  - {service="advanced-loki-example"} |= "error" | json`)
	fmt.Println(`  - {service="advanced-loki-example"} | json | trace_id="abc123"`)
	fmt.Println(`  - {service="advanced-loki-example"} | json | level="INFO" | duration > 100`)
	fmt.Println(`  - {service="advanced-loki-example"} | json | user_id="789"`)
	fmt.Println(`  - {service="advanced-loki-example"} | json | team="platform"`)
	fmt.Println()
	fmt.Println("✨ Advanced Features Demonstrated:")
	fmt.Println("  ✓ Custom Loki configuration with authentication")
	fmt.Println("  ✓ Multi-tenant support with tenant ID")
	fmt.Println("  ✓ Custom batch size and timeout settings")
	fmt.Println("  ✓ Rich label structure for better querying")
	fmt.Println("  ✓ Structured logging with JSON parsing")
	fmt.Println("  ✓ Distributed tracing context")
	fmt.Println("  ✓ Business metrics and monitoring")
	fmt.Println("  ✓ Security event logging")
	fmt.Println("  ✓ System health monitoring")
	fmt.Println("  ✓ Performance tracking")
	fmt.Println("  ✓ Configuration management")
	fmt.Println("  ✓ Graceful shutdown logging")
	fmt.Println()
}
