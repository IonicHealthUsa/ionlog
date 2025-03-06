package logcore

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

// functionData returns package, function name, file name and line number of the caller
func functionData(skip int) (pkg, function, file string, line int) {
	// Get caller information
	pc, fullFilePath, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Fprint(os.Stderr, "Failed to get caller information\n")
		return "", "", "", 0
	}

	// Get function name
	funcObj := runtime.FuncForPC(pc)
	if funcObj == nil {
		fmt.Fprint(os.Stderr, "Failed to get function object\n")
		return "", "", "", 0
	}

	fullFuncName := funcObj.Name()

	// Extract package name
	lastDotIndex := strings.LastIndexByte(fullFuncName, '.')
	if lastDotIndex < 0 {
		pkg = ""
		function = fullFuncName
	} else {
		pkg = fullFuncName[:lastDotIndex]
		function = fullFuncName[lastDotIndex+1:]

		// Get just the last part of the package path
		if lastSlashIndex := strings.LastIndexByte(pkg, '/'); lastSlashIndex >= 0 {
			pkg = pkg[lastSlashIndex+1:]
		}
	}

	// Extract just the file name from the full path
	if lastSlashIndex := strings.LastIndexByte(fullFilePath, '/'); lastSlashIndex >= 0 {
		file = fullFilePath[lastSlashIndex+1:]
	} else {
		file = fullFilePath
	}

	return pkg, function, file, line
}

func GetRecordInformation() []any {
	pkg, function, file, line := functionData(3)
	recInf := make([]any, 4)
	recInf[0] = slog.String("package", pkg)
	recInf[1] = slog.String("function", function)
	recInf[2] = slog.String("file", file)
	recInf[3] = slog.Int("line", line)
	return recInf
}
