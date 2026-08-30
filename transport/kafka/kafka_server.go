package kafka

import (
	"crypto/tls"

	"github.com/go-kratos/kratos/v2/log"
	kafkaGo "github.com/segmentio/kafka-go"

	tlsUtils "github.com/tx7do/go-utils/tls"
	broker "github.com/tx7do/kratos-transport/broker"
	kafkaBroker "github.com/tx7do/kratos-transport/broker/kafka"
	kafkaTransport "github.com/tx7do/kratos-transport/transport/kafka"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewKafkaServer creates a new Kafka server.
func NewKafkaServer(cfg *conf.Server_Kafka, opts ...kafkaTransport.ServerOption) *kafkaTransport.Server {
	if cfg == nil {
		return nil
	}

	var o []kafkaTransport.ServerOption

	if len(cfg.GetEndpoints()) != 0 {
		o = append(o, kafkaTransport.WithAddress(cfg.GetEndpoints()))
	}

	if cfg.GetCodec() != "" {
		o = append(o, kafkaTransport.WithCodec(cfg.GetCodec()))
	}

	switch cfg.AuthMechanism.(type) {
	case *conf.Server_Kafka_Plain:
		plain := cfg.GetPlain()
		o = append(o, kafkaTransport.WithPlainMechanism(
			plain.GetUsername(), plain.GetPassword(),
		))

	case *conf.Server_Kafka_Scram:
		scram := cfg.GetScram()
		o = append(o, kafkaTransport.WithScramMechanism(
			scram.GetAlgorithm(),
			scram.GetUsername(),
			scram.GetPassword(),
		))
	}

	var brokerOpts []broker.Option

	if cfg.MaxAttempts != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithMaxAttempts(int(cfg.GetMaxAttempts())))
	}
	if cfg.BatchSize != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithBatchSize(int(cfg.GetBatchSize())))
	}
	if cfg.BatchBytes != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithBatchBytes(cfg.GetBatchBytes()))
	}
	if cfg.PublishMaxAttempts != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithPublishMaxAttempts(int(cfg.GetPublishMaxAttempts())))
	}
	if cfg.BatchTimeout != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithBatchTimeout(cfg.GetBatchTimeout().AsDuration()))
	}
	if cfg.ReadTimeout != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithReadTimeout(cfg.GetReadTimeout().AsDuration()))
	}
	if cfg.WriteTimeout != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithWriteTimeout(cfg.GetWriteTimeout().AsDuration()))
	}
	if cfg.EnableOneTopicOneWriter != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithEnableOneTopicOneWriter(cfg.GetEnableOneTopicOneWriter()))
	}
	if cfg.AllowPublishAutoTopicCreation != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithAllowPublishAutoTopicCreation(cfg.GetAllowPublishAutoTopicCreation()))
	}
	if cfg.Async != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithAsync(cfg.GetAsync()))
	}
	if cfg.EnableLogger != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithEnableLogger(cfg.GetEnableLogger()))
	}
	if cfg.EnableErrorLogger != nil {
		brokerOpts = append(brokerOpts, kafkaBroker.WithEnableErrorLogger(cfg.GetEnableErrorLogger()))
	}

	if cfg.GetStartOffset() != "" {
		var rc kafkaGo.ReaderConfig
		switch cfg.GetStartOffset() {
		case "first":
			rc.StartOffset = kafkaGo.FirstOffset
		case "last":
			rc.StartOffset = kafkaGo.LastOffset
		default:
			log.Warnf("kafka: unknown start_offset %q, ignored", cfg.GetStartOffset())
		}
		brokerOpts = append(brokerOpts, kafkaBroker.WithReaderConfig(rc))
	}

	var wc kafkaBroker.WriterConfig
	wcChanged := false
	if balancer := buildBalancer(cfg.GetBalancer()); balancer != nil {
		wc.Balancer = balancer
		wcChanged = true
	}
	switch cfg.GetRequiredAcks() {
	case "one":
		wc.RequiredAcks = kafkaGo.RequireOne
		wcChanged = true
	case "all":
		wc.RequiredAcks = kafkaGo.RequireAll
		wcChanged = true
	}
	if wcChanged {
		brokerOpts = append(brokerOpts, kafkaBroker.WithWriterConfig(wc))
	}

	if len(brokerOpts) > 0 {
		o = append(o, kafkaTransport.WithBrokerOptions(brokerOpts...))
	}

	if cfg.Tls != nil {
		tlsCfg, err := loadClientTlsConfig(cfg.Tls)
		if err != nil {
			log.Errorf("kafka: load tls config failed: %v", err)
		} else {
			o = append(o, kafkaTransport.WithTLSConfig(tlsCfg))
		}
	}

	if opts != nil {
		o = append(o, opts...)
	}

	srv := kafkaTransport.NewServer(o...)

	return srv
}

func buildBalancer(name string) kafkaGo.Balancer {
	switch name {
	case "round_robin":
		return &kafkaGo.RoundRobin{}
	case "least_bytes":
		return &kafkaGo.LeastBytes{}
	case "hash":
		return &kafkaGo.Hash{}
	case "crc32":
		return &kafkaGo.CRC32Balancer{}
	case "murmur2":
		return &kafkaGo.Murmur2Balancer{}
	case "reference_hash":
		return &kafkaGo.ReferenceHash{}
	}
	return nil
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

	return tlsCfg, err
}
