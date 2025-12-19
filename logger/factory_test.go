package logger

import (
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestRegisterAndNewLogger(t *testing.T) {
	// ensure clean state: register a mock factory and then unregister at the end
	typ := Type("mock")
	_ = Register(typ, func(cfg *conf.Logger) (log.Logger, error) {
		return NewStdLogger(), nil
	})
	defer Unregister(typ)

	// ListFactories should contain our mock type
	found := false
	for _, tt := range ListFactories() {
		if tt == typ {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected ListFactories to contain %s", typ)
	}

	// NewLogger using cfg.Type should create logger
	cfg := &conf.Logger{Type: string(typ)}
	lg, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger returned error: %v", err)
	}
	if lg == nil {
		t.Fatalf("expected logger instance, got nil")
	}
}

func TestUnregisterAndUnsupported(t *testing.T) {
	typ := Type("temp")
	_ = Register(typ, func(cfg *conf.Logger) (log.Logger, error) { return NewStdLogger(), nil })
	// Unregister and verify NewLogger returns unsupported
	ok := Unregister(typ)
	if !ok {
		t.Fatalf("expected Unregister to return true")
	}
	cfg := &conf.Logger{Type: string(typ)}
	_, err := NewLogger(cfg)
	if err == nil || !strings.Contains(err.Error(), "unsupported logger type") {
		t.Fatalf("expected unsupported logger type error, got %v", err)
	}
}

func TestFactoryReturnsNilLogger(t *testing.T) {
	typ := Type("badfactory")
	_ = Register(typ, func(cfg *conf.Logger) (log.Logger, error) { return nil, nil })
	defer Unregister(typ)

	cfg := &conf.Logger{Type: string(typ)}
	_, err := NewLogger(cfg)
	if err == nil || !strings.Contains(err.Error(), "returned nil logger") {
		t.Fatalf("expected nil logger error, got %v", err)
	}
}

func TestNewLoggerProviderFallback(t *testing.T) {
	// nil cfg should return a non-nil std logger
	lg := NewLoggerProvider(nil, nil)
	if lg == nil {
		t.Fatalf("expected non-nil logger from NewLoggerProvider(nil,nil)")
	}

	// empty type cfg should also fallback to std
	cfg := &conf.Logger{}
	lg2 := NewLoggerProvider(cfg, nil)
	if lg2 == nil {
		t.Fatalf("expected non-nil logger from NewLoggerProvider(empty cfg)")
	}
}
