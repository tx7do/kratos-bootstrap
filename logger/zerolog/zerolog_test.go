package zerolog

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
)

type testWriteSyncer struct {
	output []string
}

func (x *testWriteSyncer) Write(p []byte) (n int, err error) {
	x.output = append(x.output, string(p))
	return len(p), nil
}

func TestLogger(t *testing.T) {
	syncer := &testWriteSyncer{}
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zlogger := zerolog.New(syncer)
	logger := NewZerologLogger(&zlogger)

	logger.Debug(context.Background(), "", "log", "debug")
	logger.Info(context.Background(), "", "log", "info")
	logger.Warn(context.Background(), "", "log", "warn")
	logger.Error(context.Background(), "", "log", "error")

	except := []string{
		"{\"level\":\"debug\",\"log\":\"debug\"}\n",
		"{\"level\":\"info\",\"log\":\"info\"}\n",
		"{\"level\":\"warn\",\"log\":\"warn\"}\n",
		"{\"level\":\"error\",\"log\":\"error\"}\n",
	}
	for i, s := range except {
		if s != syncer.output[i] {
			t.Logf("except=%s, got=%s", s, syncer.output[i])
			t.Fail()
		}
	}
}
