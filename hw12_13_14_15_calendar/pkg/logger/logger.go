package logger

import (
	"errors"
)

var ErrorUnknownLevel = errors.New("unknown log level")
var ErrorOpenLogFile = errors.New("error open log file")

type Level int

const (
	LevelNone Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelOut
)

// Logger интерфейс логгера
type Logger interface {
	Info(data map[string]interface{}, msg string, args ...interface{})
	Warn(data map[string]interface{}, msg string, args ...interface{})
	Error(data map[string]interface{}, msg string, args ...interface{})
	Fatal(data map[string]interface{}, msg string, args ...interface{})
	Debug(data map[string]interface{}, msg string, args ...interface{})
}

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelDebug:
		return "debug"
	}
	return "none"
}

func (l Level) Valid() bool {
	return l > LevelNone && l < LevelOut
}

func ParseLevel(level string) (Level, error) {
	var l Level
	switch level {
	case "info":
		l = LevelInfo
	case "warn":
		l = LevelWarn
	case "error":
		l = LevelError
	case "fatal":
		l = LevelFatal
	case "debug":
		l = LevelDebug
	}
	if l == LevelNone {
		return LevelNone, ErrorUnknownLevel
	}
	return l, nil
}

// Config для всех логгеров
type Config struct {
	Level     Level
	FileName  string
	IsTesting bool
}
