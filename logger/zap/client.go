package zap

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

func init() {
	_ = bLogger.Register(bLogger.Zap, func(cfg *conf.Logger) (bLogger.Logger, error) {
		return NewLogger(cfg)
	})
}

// NewLogger 创建一个新的日志记录器 - Zap
func NewLogger(cfg *conf.Logger) (bLogger.Logger, error) {
	if cfg == nil || cfg.Zap == nil {
		return nil, nil
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var encoder zapcore.Encoder
	switch strings.ToLower(cfg.Zap.GetEncoder()) {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var writeSyncer zapcore.WriteSyncer
	switch strings.ToLower(cfg.Zap.GetWriter()) {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	default:
		lumberJackLogger := &lumberjack.Logger{
			Filename:   cfg.Zap.Filename,
			MaxSize:    int(cfg.Zap.MaxSize),
			MaxBackups: int(cfg.Zap.MaxBackups),
			MaxAge:     int(cfg.Zap.MaxAge),
			Compress:   cfg.Zap.GetCompress(),
		}
		writeSyncer = zapcore.AddSync(lumberJackLogger)
	}

	var lvl = new(zapcore.Level)
	if err := lvl.UnmarshalText([]byte(cfg.Zap.Level)); err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, writeSyncer, lvl)
	l := zap.New(core).WithOptions()

	wrapped := NewZapLogger(l)

	return wrapped, nil
}
