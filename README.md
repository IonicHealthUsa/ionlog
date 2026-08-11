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

	ionlog.Configure(
		ionlog.WithStaticFields(appInfo),
		ionlog.WithWriters(ionlog.CustomOutput(os.Stdout)),
		ionlog.WithCallerInfoDepth(2), // default is 2
	)

	ionlog.Start()
	defer ionlog.Stop()

	// These logs are async
	ionlog.Infof("Test version: %v", appInfo["version"])
	ionlog.Debugf("This is a debug message: %v", "some debug info")
	ionlog.Warnf("This is a warning message: %v", "some warning info")
	ionlog.Errorf("This is an error message: %v", "some error info")

	// optional: you can turn on trace logging
	ionlog.Configure(ionlog.WithTraceMode(true))

	// Trace is a sync log
	ionlog.Tracef("This is a trace message: %v", "some trace info")
}
```

## Example Output

Every log entry is built internally as a raw JSON line. What you see on screen depends on which writer you attach.

### Without `CustomOutput` (e.g. `ionlog.WithWriters(ionlog.DefaultOutput)`)

The raw JSON line is written as-is:

```json
{"time":"2024-12-06T20:59:47.252944832-03:00","level":"INFO","msg":"Test version: 1.0.0","app":"Basic Usage","version":"1.0.0","env":"test","package":"main","function":"main","file":"main.go","line":33}

{"time":"2024-12-06T20:59:47.253012345-03:00","level":"DEBUG","msg":"This is a debug message: some debug info","app":"Basic Usage","version":"1.0.0","env":"test","package":"main","function":"main","file":"main.go","line":34}

{"time":"2024-12-06T20:59:47.253078901-03:00","level":"WARN","msg":"This is a warning message: some warning info","app":"Basic Usage","version":"1.0.0","env":"test","package":"main","function":"main","file":"main.go","line":35}

{"time":"2024-12-06T20:59:47.253145678-03:00","level":"ERROR","msg":"This is an error message: some error info","app":"Basic Usage","version":"1.0.0","env":"test","package":"main","function":"main","file":"main.go","line":36}

{"time":"2024-12-06T20:59:47.253212345-03:00","level":"TRACE","msg":"This is a trace message: some trace info","app":"Basic Usage","version":"1.0.0","env":"test","package":"main","function":"main","file":"main.go","line":43}
```

### With `CustomOutput` (e.g. `ionlog.WithWriters(ionlog.CustomOutput(os.Stdout))`)

`CustomOutput` parses that same JSON line and reformats it into a colorized, human-readable line for the terminal (colors are ANSI, not visible in this markdown):

```
2024-12-06T20:59:47-03:00 INFO [main main] Test version: 1.0.0 (main.go:33) app:Basic Usage version:1.0.0 env:test

2024-12-06T20:59:47-03:00 DEBUG [main main] This is a debug message: some debug info (main.go:34) app:Basic Usage version:1.0.0 env:test

2024-12-06T20:59:47-03:00 WARN [main main] This is a warning message: some warning info (main.go:35) app:Basic Usage version:1.0.0 env:test

2024-12-06T20:59:47-03:00 ERROR [main main] This is an error message: some error info (main.go:36) app:Basic Usage version:1.0.0 env:test

2024-12-06T20:59:47-03:00 TRACE [main main] This is a trace message: some trace info (main.go:43) app:Basic Usage version:1.0.0 env:test
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

	ionlog.Configure(
		ionlog.WithStaticFields(appInfo),
		ionlog.WithWriters(ionlog.DefaultOutput),
		// ionlog.WithLogFileRotation(ionlog.DefaultLogFolder, 1*ionlog.Mebibyte, ionlog.Daily),
		ionlog.WithQueueSize(10),
		ionlog.WithCallerInfoDepth(2), // default is 2
	)

	ionlog.Start()
	defer ionlog.Stop()

	// These logs are async
	ionlog.Infof("Test version: %v", appInfo["version"])
	ionlog.Debugf("This is a debug message: %v", "some debug info")
	ionlog.Warnf("This is a warning message: %v", "some warning info")
	ionlog.Errorf("This is an error message: %v", "some error info")

	// optional: you can turn on trace logging
	ionlog.Configure(ionlog.WithTraceMode(true))

	// Trace is a sync log
	ionlog.Tracef("This is a trace message: %v", "some trace info")

	// Turn off trace mode
	ionlog.Configure(ionlog.WithTraceMode(false))

	// Add CustomOutput to wrtiters, this will be the colorful logging in the terminal.
	ionlog.Configure(ionlog.WithWriters(ionlog.CustomOutput(os.Stdout)))
	ionlog.Info("This is a log with color")

	ionlog.Configure(ionlog.WithoutWriters(ionlog.CustomOutput(os.Stdout)))
	ionlog.Info("This is a log without color, it will be written to the default output")

	// Add a static field
	ID := "0xABC123"
	ionlog.Configure(ionlog.WithStaticFields(map[string]string{"id": ID}))
	ionlog.Infof("This log has a static field: %s", ID)

	// Remove the static field
	ionlog.Configure(ionlog.WithoutStaticFields("id"))
	ionlog.Info("This log does not have the static field 'id' anymore")

	// Configure caller stack depth (useful when using wrapper functions)
	ionlog.Configure(ionlog.WithCallerInfoDepth(3))
	ionlog.Info("This log uses a custom caller stack depth")
}
```

# Key Features
## Configuration Options

### Add a writers: Log to multiple destinations (console, files, websockets, custom writers).
```go
ionlog.Configure(
    ionlog.WithWriters(ionlog.DefaultOutput, ionlog.CustomOutput, ...),
)
```

### Remove a writer: Remove the writer by its reference.
```go
ionlog.Configure(
    ionlog.WithoutWriters(ionlog.CustomOutput, ...),
)
```

### Static Fields: Add fixed fields to all logs (e.g., service name, environment).
```go
fields := map[string]string{"service-id": "0xcafe"}
ionlog.Configure(
    ionlog.WithStaticFields(fields),
)
```

### Static Fields: Remove the static fields.
```go
ionlog.Configure(
    ionlog.WithoutStaticFields("service-id"),
)
```

### Log Rotation: Auto-rotate logs by size and time.
```go
ionlog.Configure(
    ionlog.WithLogFileRotation("logs", 100*ionlog.Mebibyte, ionlog.Hourly),
)
```

### Report Size: sets the size pf reports queue.
```go
ionlog.Configure(
    ionlog.WithQueueSize(200),
)
```

### Trace: enable or disable the trace mode.
```go
ionlog.Configure(
    ionlog.WithTraceMode(true), // or false to disable
)
```

### Caller Stack Depth: configure how many stack frames to skip when retrieving caller information.
```go
ionlog.Configure(
    ionlog.WithCallerInfoDepth(3), // default is 2
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
