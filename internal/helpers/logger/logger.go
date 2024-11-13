package logger

import (
	"fmt"
	"strings"
)

type LogLevel int

const (
	EmergencyLevel LogLevel = iota
	CriticalLevel
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

var stringToLevel = map[string]LogLevel{
	"EMERG":     EmergencyLevel,
	"EMERGENCY": EmergencyLevel,
	"CRIT":      CriticalLevel,
	"CRITICAL":  CriticalLevel,
	"ERROR":     ErrorLevel,
	"WARN":      WarningLevel,
	"WARNING":   WarningLevel,
	"INFO":      InfoLevel,
	"DEBUG":     DebugLevel,
	"TRACE":     TraceLevel,
}

var levelToString = map[LogLevel]string{
	EmergencyLevel: "EMERG",
	CriticalLevel:  "CRIT",
	ErrorLevel:     "ERROR",
	WarningLevel:   "WARN",
	InfoLevel:      "INFO",
	DebugLevel:     "DEBUG",
	TraceLevel:     "TRACE",
}

var levelToColor = map[LogLevel]string{
	EmergencyLevel: "\x1b[41m\x1b[38;5;226m",
	CriticalLevel:  "\033[31m",
	ErrorLevel:     "\x1b[38;5;208m",
	WarningLevel:   "\x1b[38;5;226m",
	InfoLevel:      "",
	DebugLevel:     "\x1b[38;5;245m",
	TraceLevel:     "\x1b[38;5;237m",
}

var curLevel LogLevel

func LevelFromString(level string) LogLevel {
	for k, v := range stringToLevel {
		if strings.EqualFold(k, level) {
			return v
		}
	}
	return ErrorLevel
}

func SetLevel(level LogLevel) {
	curLevel = level
}

func Fatal(msg string, args ...any) {
	log(EmergencyLevel, msg, args...)
	panic(fmt.Sprintf(msg, args...))
}
func Emergency(msg string, args ...any) { log(EmergencyLevel, msg, args...) }
func Critical(msg string, args ...any)  { log(CriticalLevel, msg, args...) }
func Error(msg string, args ...any)     { log(ErrorLevel, msg, args...) }
func Warning(msg string, args ...any)   { log(WarningLevel, msg, args...) }
func Info(msg string, args ...any)      { log(InfoLevel, msg, args...) }
func Debug(msg string, args ...any)     { log(DebugLevel, msg, args...) }
func Trace(msg string, args ...any)     { log(TraceLevel, msg, args...) }

func log(level LogLevel, msg string, args ...any) {
	if level > curLevel {
		return
	}
	color := levelToColor[level]
	fmt.Printf("%s[%s] %s\033[0m\n", color, levelToString[level], fmt.Sprintf(msg, args...))
}
