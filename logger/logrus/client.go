package logrus

import (
	"github.com/sirupsen/logrus"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

func init() {
	_ = bLogger.Register(bLogger.Logrus, func(cfg *conf.Logger) (bLogger.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Logrus
func NewLogger(cfg *conf.Logger) (bLogger.Logger, error) {
	if cfg == nil || cfg.Logrus == nil {
		return nil, nil
	}

	loggerLevel, err := logrus.ParseLevel(cfg.Logrus.Level)
	if err != nil {
		loggerLevel = logrus.InfoLevel
	}

	var loggerFormatter logrus.Formatter
	switch cfg.Logrus.Formatter {
	default:
		fallthrough
	case "text":
		loggerFormatter = &logrus.TextFormatter{
			DisableColors:    cfg.Logrus.DisableColors,
			DisableTimestamp: cfg.Logrus.DisableTimestamp,
			TimestampFormat:  cfg.Logrus.TimestampFormat,
		}
	case "json":
		loggerFormatter = &logrus.JSONFormatter{
			DisableTimestamp: cfg.Logrus.DisableTimestamp,
			TimestampFormat:  cfg.Logrus.TimestampFormat,
		}
	}

	l := logrus.New()
	l.Level = loggerLevel
	l.Formatter = loggerFormatter

	wrapped := NewLogrusLogger(l)
	return wrapped, nil
}
