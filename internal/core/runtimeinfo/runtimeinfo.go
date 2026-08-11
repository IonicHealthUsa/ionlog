package runtimeinfo

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type CallerInfo struct {
	File     string
	Package  string
	Function string
	Line     int
}

func GetCallerInfo(skip int) CallerInfo {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Fprint(os.Stderr, "Failed to get caller information\n")
		return CallerInfo{}
	}

	fileLastSlashIndex := strings.LastIndexByte(file, '/')

	// Get function name
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		fmt.Fprint(os.Stderr, "Failed to get caller function information\n")
		return CallerInfo{File: file[fileLastSlashIndex+1:], Line: line}
	}
	fullFuncName := fn.Name()

	lastSlashIndex := strings.LastIndexByte(fullFuncName, '/')

	fistDotIndex := strings.IndexByte(fullFuncName[lastSlashIndex+1:], '.')
	if fistDotIndex == -1 {
		// No package-qualifying dot found (e.g. a runtime-synthesized frame);
		// report the whole name as the function and leave Package empty.
		return CallerInfo{
			File:     file[fileLastSlashIndex+1:],
			Function: fullFuncName[lastSlashIndex+1:],
			Line:     line,
		}
	}
	pkgEnd := lastSlashIndex + 1 + fistDotIndex

	return CallerInfo{
		File:     file[fileLastSlashIndex+1:],
		Package:  fullFuncName[lastSlashIndex+1 : pkgEnd],
		Function: fullFuncName[pkgEnd+1:],
		Line:     line,
	}
}
