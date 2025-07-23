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
	fullFuncName := runtime.FuncForPC(pc).Name()

	lastSlashIndex := strings.LastIndexByte(fullFuncName, '/')

	fistDotIndex := strings.IndexByte(fullFuncName[lastSlashIndex+1:], '.')
	pkgEnd := lastSlashIndex + 1 + fistDotIndex

	return CallerInfo{
		File:     file[fileLastSlashIndex+1:],
		Package:  fullFuncName[lastSlashIndex+1 : pkgEnd],
		Function: fullFuncName[pkgEnd+1:],
		Line:     line,
	}
}
