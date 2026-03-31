package kafka

import (
	"crypto/tls"

	tlsUtils "github.com/tx7do/go-utils/tls"
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

	if cfg.MaxAttempts != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithMaxAttempts(int(cfg.GetMaxAttempts())),
		)
	}

	if cfg.BatchSize != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithBatchSize(int(cfg.GetBatchSize())),
		)
	}

	if cfg.BatchBytes != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithBatchBytes(cfg.GetBatchBytes()),
		)
	}

	if cfg.PublishMaxAttempts != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithPublishMaxAttempts(int(cfg.GetPublishMaxAttempts())),
		)
	}

	if cfg.BatchTimeout != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithBatchTimeout(cfg.GetBatchTimeout().AsDuration()),
		)
	}
	if cfg.ReadTimeout != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithReadTimeout(cfg.GetReadTimeout().AsDuration()),
		)
	}
	if cfg.WriteTimeout != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithWriteTimeout(cfg.GetWriteTimeout().AsDuration()),
		)
	}

	if cfg.EnableOneTopicOneWriter != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithEnableOneTopicOneWriter(cfg.GetEnableOneTopicOneWriter()),
		)
	}

	if cfg.AllowPublishAutoTopicCreation != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithAllowPublishAutoTopicCreation(cfg.GetAllowPublishAutoTopicCreation()),
		)
	}

	if cfg.Async != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithAsync(cfg.GetAsync()),
		)
	}

	if cfg.EnableLogger != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithEnableLogger(cfg.GetEnableLogger()),
		)
	}

	if cfg.EnableErrorLogger != nil {
		kafkaTransport.WithBrokerOptions(
			kafkaBroker.WithEnableErrorLogger(cfg.GetEnableErrorLogger()),
		)
	}

	if cfg.Tls != nil {
		if tlsCfg, err := loadClientTlsConfig(cfg.Tls); err != nil {
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
