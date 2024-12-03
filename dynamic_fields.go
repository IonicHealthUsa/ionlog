package ionlog

import (
	"log/slog"
	"runtime"
	"strings"
)

func function_data(skip int) (pkg, function, file string, line int) {
	pc, file, line, _ := runtime.Caller(skip)
	data := runtime.FuncForPC(pc).Name()
	pkg = data[:strings.LastIndexByte(data, '.')]
	if strings.Contains(pkg, "/") {
		pkg = pkg[strings.LastIndexByte(pkg, '/')+1:]
	}
	function = data[strings.LastIndexByte(data, '.')+1:]
	file = file[strings.LastIndexByte(file, '/')+1:]
	return
}

func getRecordInformation() []any{
	pkg, function, file, line := function_data(3)
	recInf := make([]any, 4)
	recInf[0] = slog.String("package", pkg)
	recInf[1] = slog.String("function", function)
	recInf[2] = slog.String("file", file)
	recInf[3] = slog.Int("line", line)
	return recInf
}
