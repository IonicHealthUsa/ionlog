package main

import "ionlog"

func main() {
	ionlog.SetLogAttributes(
		ionlog.WithTargets(ionlog.Stdout() /* socket, files ... */),
		ionlog.WithAttrs(map[string]string{
			"e-KVM_id": "1234567890",
			"...":      "...",
		}),
	)

	ionlog.Info("This is a log message, with %s", "args")
	ionlog.Error("This is a log message, with %v", "args")
	ionlog.Debug("This is a log message, with %s", "args")
	ionlog.Warn("This is a log message, with %s", "args")
}
