package logcore

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func GetCallerInfo(skip int) *callerInfo {
	var (
		file     string
		pkg      string
		function string
	)

	var line int

	// Get caller information
	pc, fullFilePath, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Fprint(os.Stderr, "Failed to get caller information\n")
		return nil
	}

	// Get function name
	funcObj := runtime.FuncForPC(pc)
	if funcObj == nil {
		fmt.Fprint(os.Stderr, "Failed to get function object\n")
		return nil
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

	return &callerInfo{
		File:        file,
		PackageName: pkg,
		Function:    function,
		Line:        line,
	}
}
