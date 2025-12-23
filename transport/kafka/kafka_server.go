package kafka

import (
	"github.com/tx7do/kratos-transport/transport/kafka"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewKafkaServer creates a new Kafka server.
func NewKafkaServer(cfg *conf.Server_Kafka, opts ...kafka.ServerOption) *kafka.Server {
	if cfg == nil {
		return nil
	}

	var o []kafka.ServerOption

	if len(cfg.GetEndpoints()) != 0 {
		o = append(o, kafka.WithAddress(cfg.GetEndpoints()))
	}

	if cfg.GetCodec() != "" {
		o = append(o, kafka.WithCodec(cfg.GetCodec()))
	}

	if opts != nil {
		o = append(o, opts...)
	}

	srv := kafka.NewServer(o...)

	return srv
}
