package sse

import (
	"github.com/tx7do/kratos-transport/transport/sse"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewSseServer creates a new SSE server.
func NewSseServer(cfg *conf.Server_SSE, opts ...sse.ServerOption) *sse.Server {
	if cfg == nil {
		return nil
	}

	var o []sse.ServerOption

	if cfg.GetNetwork() != "" {
		o = append(o, sse.WithNetwork(cfg.GetNetwork()))
	}
	if cfg.GetAddr() != "" {
		o = append(o, sse.WithAddress(cfg.GetAddr()))
	}
	if cfg.GetPath() != "" {
		o = append(o, sse.WithPath(cfg.GetPath()))
	}
	if cfg.GetCodec() != "" {
		o = append(o, sse.WithCodec(cfg.GetCodec()))
	}

	if cfg.Timeout != nil {
		o = append(o, sse.WithTimeout(cfg.GetTimeout().AsDuration()))
	}
	if cfg.EventTtl != nil {
		o = append(o, sse.WithEventTTL(cfg.GetEventTtl().AsDuration()))
	}

	o = append(o,
		sse.WithAutoStream(cfg.GetAutoStream()),
		sse.WithAutoReply(cfg.GetAutoReply()),
		sse.WithSplitData(cfg.GetSplitData()),
		sse.WithEncodeBase64(cfg.GetEncodeBase64()),
	)

	if opts != nil {
		o = append(o, opts...)
	}

	srv := sse.NewServer(o...)

	return srv
}
