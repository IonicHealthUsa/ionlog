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
    ionlog.SetLogAttributes(
        ionlog.WithTargets(ionlog.DefaultOutput), // Log to console
        ionlog.WithStaticFields(map[string]string{"service": "my-app"}),
        ionlog.WithLogFileRotation("logs", 10*ionlog.Mebibyte, ionlog.Daily),
    )

    ionlog.Start()
    defer ionlog.Stop()

    ionlog.Info("Application started")
}
```

# Advanced Usage
```go
package main

import (
	"github.com/IonicHealthUsa/ionlog"
)

func main() {
	// SetAttributes set the log attributes, and other configurations
	ionlog.SetAttributes(
		// WithWriters sets the write targets for the logger, every log will be written
		// to these targets.
		ionlog.WithWriters(
			ionlog.CustomOutput,
			// a websocket
			// a file
			// your custom writer
		),

		// (Optinal) WithoutTarget remove the writer targets for the logger. Pass the pointer of the writer.
		ionlog.WithoutWriters(
			// previously writer defined
		),

		// (Optional) WithStaticFields sets the static fields for the logger, every log will have these fields.
		ionlog.WithStaticFields(map[string]string{
			"computer-id": "1234",
			// your custom fields
		}),

		// (Optional) WithoutStaticFields remove the static fields for the logger. Use the key of the static field to remove.
		ionlog.WithoutStaticFields(
			// previously static field defined
		),

		// (Optional) WithLogFileRotation enables log file rotation, specifying the directory where log files will be stored, the maximum size of the log folder in bytes, and the rotation frequency.
		// This internal log rotation system appends the log file to the specified targets and automatically rotates logs based on the provided configuration,
		// ensuring the total size of the log folder does not exceed the specified maximum (e.g., 10MB in this case).
		ionlog.WithLogFileRotation("logs", 10*ionlog.Mebibyte, ionlog.Daily),

		// (Optional) WithQueueSize sets the size of the reports queue, which stores logs before sending them to a file descriptor. For default, the size of report is 100.
		ionlog.WithQueueSize(120),

		// (Optional) WithTraceMode enables trace log mode. For default, the trace mode is disable, to enable is need pass a tru boolean.
		ionlog.WithTraceMode(true),
	)

	// Start the logger service
	ionlog.Start()

	// Stops the logger service when the main function ends
	defer ionlog.Stop()

	// output: {"time":"2024-12-06T20:59:47.252944832-03:00","level":"INFO","msg":"This log level is: info","computer-id":"1234","package":"main","function":"main","file":"main.go","line":38}
	ionlog.Infof("This log level is: %v", "info")
	ionlog.Errorf("This log level is: %v", "error")
	ionlog.Warnf("This log level is: %v", "warn")
	ionlog.Debugf("This log level is: %v", "debug")
	ionlog.Tracef("This log level is: %v", "trace")

	ionlog.Info("This log level is a simple info log")
	ionlog.Error("This log level is a simple error log")
	ionlog.Warn("This log level is a simple warn log")
	ionlog.Debug("This log level is a simple debug log")
	ionlog.Trace("This log level is a simple trace log")

	status := "NOT OK"
	for i := range 10 {
		ionlog.LogOnceInfo("Process Started!") // This will be logged only once
		ionlog.LogOnceDebugf("count: %v", i)   // Log every time i changes
		if i == 5 {
			status = "OK"
		}
		ionlog.LogOnceInfof("status: %v", status) // Log once "NOT OK", log once "OK"
	}
}
```

# Key Features
## Configuration Options

### Targets: Log to multiple destinations (console, files, websockets, custom writers).
```go
ionlog.WithTargets(ionlog.DefaultOutput, myCustomWriter)
```

### Targets: Remove the target.
```go
ionlog.WithoutTargets(ionlog.DefaultOutput, myCustomWriter)
```

### Static Fields: Add fixed fields to all logs (e.g., service name, environment).
```go
ionlog.WithStaticFields(map[string]string{"env": "production"})
```

### Static Fields: Remove the static fields.
```go
ionlog.WithStaticFields("env")
```

### Log Rotation: Auto-rotate logs by size and time.
```go
ionlog.WithLogFileRotation("logs", 100*ionlog.Mebibyte, ionlog.Hourly)
```

### Report Size: sets the size pf reports queue.
```go
ionlog.WithQueueSize(200)
```

### Trace: enable or disable the trace mode.
```go
ionlog.WithTraceMode(true)
```

## Logging Functions
- Levels: Debug, Info, Warn, Error.
```go
ionlog.Infof("User %s logged in", "Alice")
ionlog.Error("Connection failed")
```

- The trace level is optional. It is necessary to enable.
```go
ionlog.Trace("Trace the path")
```

## Structured Output: Logs are emitted as JSON with metadata:
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

### Log Once: Write a message only once during execution.
```go
ionlog.LogOnceInfo("Initialization complete")
```

### Log on Change: Only log when the value changes.
```go
status := "STARTING"
ionlog.LogOnceInfof("status: %s", status) // Logs once
ionlog.LogOnceInfof("status: %s", status) // Will not log

status = "RUNNING"
ionlog.LogOnceInfof("status: %s", status) // Logs again
ionlog.LogOnceInfof("status: %s", status) // Will not log
```

## Lifecycle Management:

- Start() initializes the logger
```go
ionlog.Start()
```

- Stop() closes the logger when the program ends
```go
ionlog.Stop()
```


# Internal Logging system:
- Internal logs are handled by the slog package, and outputed to the os.Stdout by default.


# Process Flow Diagram
```mermaid
sequenceDiagram
    participant P as User/Third-Party Program
    participant I as Ionlog Core
    participant S as Ionlog Settings
    participant H as Processing Handlers
    participant W as Writer Interface
    participant O as Output

    P->>I: 1. Initialize Ionlog (e.g., ionlog.SetDefault())
    activate I
    I->>S: 2. Load initial settings (Level, Format, Handlers, Writers, etc.)
    activate S
    S-->>I: 3. Return loaded settings
    deactivate S
    I-->>P: 4. Ionlog ready for use
    deactivate I


    P->>I: 5. Start logger service
	I->>H: 6. Inicitialize the handler


    P->>I: 7. Call log function (e.g., log.Info(), log.Error("msg"))
    activate I
    I->>S: 8. Query settings (trace mode, static fields, active writers)
    activate S
    S-->>I: 9. Return relevant settings for the event
    deactivate S

    alt Log level of event was configured 
        I->>H: 10. Forward log event to configured Handlers
        activate H
        H-->>H: 11. Format log message and metadata (e.g., JSON, plain text)
        Note over H: Other processing handlers (e.g., enrichment) may also apply here.

        opt Log Once Handler
            H-->>H: 12. Filter log event (if applicable, based on tags, fields, etc.)
            alt Event Filtered and Discarded
                H--xI: 13. Event discarded by filter
                Note over P,H: The log is ignored because it was already logged before.
                I-->>P: 14. Return (log event ignored)
            end
        end

        H-->>I: 15. Return processed log event
        deactivate H

        I->>W: 16. Send processed log event to configured Writers
        activate W
        opt Custom Output Processing
            W-->>W: 17. Customize log for the specific output
        end
        W->>O: 18. Write/Display log to concrete Output (Console, File, File Descriptor, etc.)
        activate O
        O-->>W: 19. Confirmation of write (or error)
        deactivate O
        W-->>I: 20. Confirmation of write to Core
        deactivate W
        I-->>P: 21. End of log operation
    else Log level of event was not configured 
        I-->>P: 10. Log event ignored (level below configured)
        deactivate I
    end


    opt Stop logger service
    P->>I: 22. Stop Ionlog
    activate I
    I->>H: 23. Signal Handlers to shut down (e.g., flush buffers, close resources)
    activate H
    H-->>I: 24. Handlers shut down confirmation
    deactivate H
    I->>W: 25. Signal Writers to shut down (e.g., flush pending writes, close file handles)
    activate W
    W->>O: 26. Perform final flush/cleanup on Outputs
    activate O
    O-->>W: 27. Outputs cleanup complete
    deactivate O
    W-->>I: 28. Writers shut down confirmation
    deactivate W
    I-->>P: 29. Ionlog shutdown complete
    deactivate I

    end
```
