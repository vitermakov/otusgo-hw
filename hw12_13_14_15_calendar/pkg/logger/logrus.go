package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
)

type Logrus struct {
	logger *logrus.Logger
}

func NewLogrus(cfg Config) (Logger, error) {
	inst := logrus.New()
	if !cfg.Level.Valid() {
		return nil, ErrorUnknownLevel
	}
	var level logrus.Level
	switch cfg.Level {
	case LevelNone:
		level = logrus.PanicLevel
	case LevelInfo:
		level = logrus.InfoLevel
	case LevelWarn:
		level = logrus.WarnLevel
	case LevelError:
		level = logrus.ErrorLevel
	case LevelFatal:
		level = logrus.FatalLevel
	case LevelDebug:
		level = logrus.DebugLevel
	}
	if len(cfg.FileName) > 0 {
		f, err := os.OpenFile(cfg.FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, ErrorOpenLogFile
		}
		if !cfg.IsTesting {
			inst.SetOutput(io.MultiWriter(f, os.Stdout))
		} else {
			inst.SetOutput(f)
		}
	}
	inst.SetLevel(level)

	if cfg.IsTesting {
		inst.ExitFunc = func(i int) {
			log.Printf("logrus fatal exit(%d)\n", i)
		}
	}
	return &Logrus{logger: inst}, nil
}

func (l Logrus) with(data map[string]interface{}) *logrus.Entry {
	if data == nil {
		data = logrus.Fields{}
	}
	return l.logger.WithFields(data)
}

func (l Logrus) Debug(data map[string]interface{}, msg string, args ...interface{}) {
	l.with(data).Debugf(msg, args...)
}

func (l Logrus) Info(data map[string]interface{}, msg string, args ...interface{}) {
	l.with(data).Infof(msg, args...)
}

func (l Logrus) Warn(data map[string]interface{}, msg string, args ...interface{}) {
	l.with(data).Warnf(msg, args...)
}

func (l Logrus) Error(data map[string]interface{}, msg string, args ...interface{}) {
	l.with(data).Errorf(msg, args...)
}

func (l Logrus) Fatal(data map[string]interface{}, msg string, args ...interface{}) {
	l.with(data).Fatalf(msg, args...)
}
