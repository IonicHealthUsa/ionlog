package logcore

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func GetCallerInfo(skip int) callerInfo {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Fprint(os.Stderr, "Failed to get caller information\n")
		return callerInfo{}
	}

	fileLastSlashIndex := strings.LastIndexByte(file, '/')

	// Get function name
	fullFuncName := runtime.FuncForPC(pc).Name()

	lastSlashIndex := strings.LastIndexByte(fullFuncName, '/')

	fistDotIndex := strings.IndexByte(fullFuncName[lastSlashIndex+1:], '.')
	pkgEnd := lastSlashIndex + 1 + fistDotIndex

	return callerInfo{
		File:        file[fileLastSlashIndex+1:],
		PackageName: fullFuncName[lastSlashIndex+1 : pkgEnd],
		Function:    fullFuncName[pkgEnd+1:],
		Line:        line,
	}
}
