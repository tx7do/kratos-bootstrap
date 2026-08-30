package sse

import (
	"crypto/tls"

	"github.com/go-kratos/kratos/v2/log"

	tlsUtils "github.com/tx7do/go-utils/tls"
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

	if cfg.BufferSize != nil && cfg.GetBufferSize() > 0 {
		o = append(o, sse.WithBufferSize(int(cfg.GetBufferSize())))
	}
	if n := len(cfg.GetHeaders()); n > 0 {
		headers := make(map[string]string, n)
		for k, v := range cfg.GetHeaders() {
			headers[k] = v
		}
		o = append(o, sse.WithHeaders(headers))
	}

	if cfg.Tls != nil {
		tlsCfg, err := loadClientTlsConfig(cfg.Tls)
		if err != nil {
			log.Errorf("sse: load tls config failed: %v", err)
		} else {
			o = append(o, sse.WithTLSConfig(tlsCfg))
		}
	}

	if opts != nil {
		o = append(o, opts...)
	}

	srv := sse.NewServer(o...)

	return srv
}

func loadClientTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var tlsCfg *tls.Config
	var err error

	if cfg.File != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigFile(
			cfg.File.GetKeyPath(),
			cfg.File.GetCertPath(),
			cfg.File.GetCaPath(),
		); err != nil {
			return nil, err
		}
	} else if cfg.Config != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigString(
			cfg.Config.GetKeyPem(),
			cfg.Config.GetCertPem(),
			cfg.Config.GetCaPem(),
		); err != nil {
			return nil, err
		}
	}

	if tlsCfg != nil && cfg.GetInsecureSkipVerify() {
		tlsCfg.InsecureSkipVerify = true
	}

	return tlsCfg, nil
}
