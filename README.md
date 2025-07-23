# ionlog

A flexible and structured logging library for Go with dynamic controls.

## Installation

```bash
go get github.com/IonicHealthUsa/ionlog
```

# Basic Usage
```go
package main

import "github.com/IonicHealthUsa/ionlog"

func main() {
	appInfo := map[string]string{
		"app":     "Basic Usage",
		"version": "1.0.0",
		"env":     "test",
	}

	ionlog.SetAttributes(
		ionlog.WithStaticFields(appInfo),
		ionlog.WithWriters(ionlog.CustomOutput),
	)

	ionlog.Start()
	defer ionlog.Stop()

	// These logs are async
	ionlog.Infof("Test version: %v", appInfo["version"])
	ionlog.Debugf("This is a debug message: %v", "some debug info")
	ionlog.Warnf("This is a warning message: %v", "some warning info")
	ionlog.Errorf("This is an error message: %v", "some error info")

	// optional: you can turn on trace logging
	ionlog.SetAttributes(ionlog.WithTraceMode(true))

	// Trace is a sync log
	ionlog.Tracef("This is a trace message: %v", "some trace info")
}
```

# Advanced Usage
```go
package main

import "github.com/IonicHealthUsa/ionlog"

func main() {
	appInfo := map[string]string{
		"app":     "Basic Usage",
		"version": "1.0.0",
		"env":     "test",
	}

	ionlog.SetAttributes(
		ionlog.WithStaticFields(appInfo),
		ionlog.WithWriters(ionlog.DefaultOutput),
		// ionlog.WithLogFileRotation(ionlog.DefaultLogFolder, 1*ionlog.Mebibyte, ionlog.Daily),
		ionlog.WithQueueSize(10),
	)

	ionlog.Start()
	defer ionlog.Stop()

	// These logs are async
	ionlog.Infof("Test version: %v", appInfo["version"])
	ionlog.Debugf("This is a debug message: %v", "some debug info")
	ionlog.Warnf("This is a warning message: %v", "some warning info")
	ionlog.Errorf("This is an error message: %v", "some error info")

	// optional: you can turn on trace logging
	ionlog.SetAttributes(ionlog.WithTraceMode(true))

	// Trace is a sync log
	ionlog.Tracef("This is a trace message: %v", "some trace info")

	// Turn off trace mode
	ionlog.SetAttributes(ionlog.WithTraceMode(false))

	// Add CustomOutput to wrtiters, this will be the colorful logging in the terminal.
	ionlog.SetAttributes(ionlog.WithWriters(ionlog.CustomOutput))
	ionlog.Info("This is a log with color")

	ionlog.SetAttributes(ionlog.WithoutWriters(ionlog.CustomOutput))
	ionlog.Info("This is a log without color, it will be written to the default output")

	// Add a static field
	ID := "0xABC123"
	ionlog.SetAttributes(ionlog.WithStaticFields(map[string]string{"id": ID}))
	ionlog.Infof("This log has a static field: %s", ID)

	// Remove the static field
	ionlog.SetAttributes(ionlog.WithoutStaticFields("id"))
	ionlog.Info("This log does not have the static field 'id' anymore")
}
```

# Key Features
## Configuration Options

### Add a writers: Log to multiple destinations (console, files, websockets, custom writers).
```go
ionlog.SetAttributes(
    ionlog.WithWriters(ionlog.DefaultOutput, ionlog.CustomOutput, ...),
)
```

### Remove a writer: Remove the writer by its reference.
```go
ionlog.SetAttributes(
    ionlog.WithoutWriters(ionlog.CustomOutput, ...),
)
```

### Static Fields: Add fixed fields to all logs (e.g., service name, environment).
```go
fields := map[string]string{"service-id": "0xcafe"}
ionlog.SetAttributes(
    ionlog.WithStaticFields(fields),
)
```

### Static Fields: Remove the static fields.
```go
ionlog.SetAttributes(
    ionlog.WithoutStaticFields("service-id"),
)
```

### Log Rotation: Auto-rotate logs by size and time.
```go
ionlog.SetAttributes(
    ionlog.WithLogFileRotation("logs", 100*ionlog.Mebibyte, ionlog.Hourly),
)
```

### Report Size: sets the size pf reports queue.
```go
ionlog.SetAttributes(
    ionlog.WithQueueSize(200),
)
```

### Trace: enable or disable the trace mode.
```go
ionlog.SetAttributes(
    ionlog.WithTraceMode(true), // or false to disable
)
```

## Logging Functions
- Levels: Debug, Info, Warn, Error.
```go
ionlog.Debug("Debugging information")
ionlog.Infof("User %s logged in", "Alice")
ionlog.Warn("Low disk space warning")
ionlog.Error("Connection failed")
```

- The trace level is optional. It is necessary to enable.
```go
ionlog.Trace("Trace the path")
```

## Structured Output: Logs are emitted as JSON with metadata ("serivce-id" is an example of static fields):
```json
{
	"time":"2024-12-06T20:59:47.252944832-03:00",
	"level":"INFO",
	"msg": "User Alice logged in",
	"service-id":"0xcafe",
	"package":"main",
	"function":"main",
	"file":"main.go",
	"line":42
}
```

## Special Logging

### Log Once: Write a message only once during execution (levels: Debug, Info, Warn, Error).
```go
ionlog.LogOnceInfo("Initialization complete")
```

## Lifecycle Management:

- Start() initializes the logger
```go
ionlog.Start()
```

- Stop() ends the logger service, flushing any pending logs and reset the log instance.
```go
ionlog.Stop()
```

# Process Flow Diagram
TODO
<!-- ```mermaid -->
<!---->
<!-- ``` -->
