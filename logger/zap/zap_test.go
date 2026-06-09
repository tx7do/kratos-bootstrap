package zap

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testWriteSyncer struct {
	output []string
}

func (x *testWriteSyncer) Write(p []byte) (n int, err error) {
	x.output = append(x.output, string(p))
	return len(p), nil
}

func (x *testWriteSyncer) Sync() error {
	return nil
}

func TestLogger(t *testing.T) {
	syncer := &testWriteSyncer{}
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), syncer, zap.DebugLevel)
	zlogger := zap.New(core).WithOptions()
	logger := NewZapLogger(zlogger)

	defer func() { _ = logger.Sync() }()

	logger.Debug(context.Background(), "", "log", "debug")
	logger.Info(context.Background(), "", "log", "info")
	logger.Warn(context.Background(), "", "log", "warn")
	logger.Error(context.Background(), "", "log", "error")
	logger.Info(context.Background(), "hello world")

	except := []string{
		"{\"level\":\"debug\",\"msg\":\"\",\"log\":\"debug\"}\n",
		"{\"level\":\"info\",\"msg\":\"\",\"log\":\"info\"}\n",
		"{\"level\":\"warn\",\"msg\":\"\",\"log\":\"warn\"}\n",
		"{\"level\":\"error\",\"msg\":\"\",\"log\":\"error\"}\n",
		"{\"level\":\"info\",\"msg\":\"hello world\"}\n",
	}
	for i, s := range except {
		if s != syncer.output[i] {
			t.Logf("except=%s, got=%s", s, syncer.output[i])
			t.Fail()
		}
	}
}
