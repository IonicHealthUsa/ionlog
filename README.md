# IonLog

# Usage
```go
package main

import (
	"github.com/IonicHealthUsa/ionlog/pkg/ionlog"
)

func main() {
	// Set the log attributes, and other configurations
	ionlog.SetLogAttributes(
		// WithTargets sets the write targets for the logger, every log will be written
		// to these targets.
		ionlog.WithTargets(
			ionlog.DefaultOutput(),
			// a websocket
			// a file
			// your custom writer
		),

		// (Optional) WithStaicFields sets the static fields for the logger, every log will have these fields.
		ionlog.WithStaicFields(map[string]string{
			"computer-id": "1234",
			// your custom fields
		}),

		// (Optional) WithLogFileRotation sets the log file rotation period and the folder where the log files will be stored.
		// This is a internal log file rotation system, when optionally used, it will append the log file to the targets, and
		// will rotate it automatically.
		ionlog.WithLogFileRotation("logs", ionlog.Daily),
	)

	// Start the logger service
	ionlog.Start()

	// Stops the logger service when the main function ends
	defer ionlog.Stop()

	// output: {"time":"2024-12-06T20:59:47.252944832-03:00","level":"INFO","msg":"This log level is: info","computer-id":"1234","package":"main","function":"main","file":"main.go","line":38}
	ionlog.Info("This log level is: %v", "info")
	ionlog.Error("This log level is: %v", "error")
	ionlog.Warn("This log level is: %v", "warn")
	ionlog.Debug("This log level is: %v", "debug")

	status := "NOT OK"
	for i := 0; i < 10; i++ {
		ionlog.LogOnceInfo("Process Started!")  // This will be logged only once
		ionlog.LogOnChangeDebug("count: %v", i) // Log every time i changes
		if i == 5 {
			status = "OK"
		}
		ionlog.LogOnChangeInfo("status: %v", status) // Log once "NOT OK", log once "OK"
	}
}
```
# Library Import and Configuration:
The library is imported from github.com/IonicHealthUsa/ionlog/pkg/ionlog.  
Configuration is done using SetLogAttributes() method with several options:

## Log Targets:
ionlog.WithTargets() allows setting multiple log output destinations.  
It supports:

- Default output
- Websocket
- File writing
- Custom writers

## Static Fields:
WithStaticFields() adds consistent metadata to all log entries.  
In this example, a "computer-id" field is added to every log

## Log File Rotation:
WithLogFileRotation() configures automatic log file management.
- Sets log file storage location to "logs" folder
- Rotation period set to daily

# Logging Methods:
Standard log levels: Info(), Error(), Warn(), Debug();  
Special logging methods:  

Logs a message only once:
- LogOnceInfo()
- LogOnceError()
- LogOnceWarn()
- LogOnceDebug()

Logs when the message changes:
- LogOnChangeInfo()
- LogOnChangeError()
- LogOnChangeWarn()
- LogOnChangeDebug()

# Lifecycle Management:
- Start() initializes the logger
- Stop() (deferred) closes the logger when the program ends

# Log Format:
- Produces JSON-formatted logs
- Includes timestamp, log level, message
- Adds static fields, package, function, file, and line information

# Internal Logging system:
- Internal logs are handled by the slog package, and outputed to the os.Stdout by default.
