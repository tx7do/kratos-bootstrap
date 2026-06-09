package logrus

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	bLogger "github.com/tx7do/kratos-bootstrap/logger"
)

func TestLoggerLog(t *testing.T) {
	tests := map[string]struct {
		loggerLevel logrus.Level
		call        func(l bLogger.Logger)
		want        string
	}{
		"info with fields": {
			loggerLevel: logrus.InfoLevel,
			call: func(l bLogger.Logger) {
				l.Info(nil, "1", "case", "json format")
			},
			want: `{"case":"json format","level":"info","msg":"1"`,
		},
		"level unmatch": {
			loggerLevel: logrus.InfoLevel,
			call: func(l bLogger.Logger) {
				l.Debug(nil, "1", "case", "level unmatch")
			},
			want: "",
		},
		"no tags": {
			loggerLevel: logrus.InfoLevel,
			call: func(l bLogger.Logger) {
				l.Info(nil, "1")
			},
			want: `{"level":"info","msg":"1"`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			output := new(bytes.Buffer)
			logger := logrus.New()
			logger.Level = test.loggerLevel
			logger.Out = output
			logger.Formatter = &logrus.JSONFormatter{}
			wrapped := NewLogrusLogger(logger)
			test.call(wrapped)

			if !strings.HasPrefix(output.String(), test.want) {
				t.Errorf("TestName(%s): %q has not prefix %q", name, output.String(), test.want)
			}
		})
	}
}
