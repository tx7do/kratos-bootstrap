package zerolog

import (
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/rs/zerolog"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/logger"
)

func init() {
	_ = logger.Register(logger.Zerelog, func(cfg *conf.Logger) (log.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Zerolog
func NewLogger(cfg *conf.Logger) (log.Logger, error) {
	if cfg == nil || cfg.Zerolog == nil {
		return nil, nil
	}

	// 根据配置设置全局级别（可选）
	if lvl := cfg.Zerolog.Level; lvl != "" {
		switch strings.ToLower(lvl) {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn", "warning":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		default:
			// 未识别则保持默认
		}
	}

	if cfg.Zerolog.TimeFieldFormat != "" {
		zerolog.TimeFieldFormat = cfg.Zerolog.TimeFieldFormat
	}
	if cfg.Zerolog.TimestampFieldName != "" {
		zerolog.TimestampFieldName = cfg.Zerolog.TimestampFieldName
	}
	if cfg.Zerolog.LevelFieldName != "" {
		zerolog.LevelFieldName = cfg.Zerolog.LevelFieldName
	}
	if cfg.Zerolog.MessageFieldName != "" {
		zerolog.MessageFieldName = cfg.Zerolog.MessageFieldName
	}

	w, err := NewWriter(strings.ToLower(cfg.Zerolog.GetWriter()), cfg.Zerolog.GetFilename(), nil)
	if err != nil {
		return nil, err
	}

	// 创建基础 zerolog.Logger
	z := zerolog.New(w).With().Timestamp().Logger()

	wrapped := NewZerologLogger(&z)

	return wrapped, nil
}
