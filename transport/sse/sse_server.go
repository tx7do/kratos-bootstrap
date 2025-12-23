package asynq

import (
	"github.com/tx7do/kratos-transport/transport/sse"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewSseServer creates a new SSE server.
func NewSseServer(cfg *conf.Server, opts ...sse.ServerOption) *sse.Server {
	if cfg == nil || cfg.Asynq == nil {
		return nil
	}

	var o []sse.ServerOption

	if cfg.Sse.GetNetwork() != "" {
		o = append(o, sse.WithNetwork(cfg.Sse.GetNetwork()))
	}
	if cfg.Sse.GetAddr() != "" {
		o = append(o, sse.WithAddress(cfg.Sse.GetAddr()))
	}
	if cfg.Sse.GetPath() != "" {
		o = append(o, sse.WithPath(cfg.Sse.GetPath()))
	}
	if cfg.Sse.GetCodec() != "" {
		o = append(o, sse.WithCodec(cfg.Sse.GetCodec()))
	}

	if cfg.Sse.Timeout != nil {
		o = append(o, sse.WithTimeout(cfg.Sse.GetTimeout().AsDuration()))
	}
	if cfg.Sse.EventTtl != nil {
		o = append(o, sse.WithEventTTL(cfg.Sse.GetEventTtl().AsDuration()))
	}

	o = append(o,
		sse.WithAutoStream(cfg.Sse.GetAutoStream()),
		sse.WithAutoReply(cfg.Sse.GetAutoReply()),
		sse.WithSplitData(cfg.Sse.GetSplitData()),
		sse.WithEncodeBase64(cfg.Sse.GetEncodeBase64()),
	)

	if opts != nil {
		o = append(o, opts...)
	}

	srv := sse.NewServer(o...)

	return srv
}
