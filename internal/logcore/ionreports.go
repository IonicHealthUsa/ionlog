package logcore

import "time"

type callerInfo struct {
	File        string
	PackageName string
	Function    string
	Line        int
}

type IonReport struct {
	Datetime time.Time
	Level    Level
	Msg      string
	callerInfo
}
