package asynq

import (
	"github.com/tx7do/kratos-transport/transport/kafka"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewKafkaServer creates a new Kafka server.
func NewKafkaServer(cfg *conf.Server, opts ...kafka.ServerOption) *kafka.Server {
	if cfg == nil || cfg.Kafka == nil {
		return nil
	}

	var o []kafka.ServerOption

	if len(cfg.Kafka.GetEndpoints()) != 0 {
		o = append(o, kafka.WithAddress(cfg.Kafka.GetEndpoints()))
	}

	if cfg.Kafka.GetCodec() != "" {
		o = append(o, kafka.WithCodec(cfg.Kafka.GetCodec()))
	}

	if opts != nil {
		o = append(o, opts...)
	}

	srv := kafka.NewServer(o...)

	return srv
}
