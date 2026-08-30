package mqtt

import (
	"crypto/tls"

	"github.com/go-kratos/kratos/v2/log"
	tlsUtils "github.com/tx7do/go-utils/tls"
	broker "github.com/tx7do/kratos-transport/broker"

	mqttBroker "github.com/tx7do/kratos-transport/broker/mqtt"
	"github.com/tx7do/kratos-transport/transport/mqtt"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewMqttServer creates a new MQTT server.
func NewMqttServer(cfg *conf.Server_Mqtt, opts ...mqtt.ServerOption) *mqtt.Server {
	if cfg == nil {
		return nil
	}

	var o []mqtt.ServerOption

	if cfg.GetCodec() != "" {
		o = append(o, mqtt.WithCodec(cfg.GetCodec()))
	}

	if cfg.GetEndpoint() != "" {
		o = append(o, mqtt.WithAddress([]string{cfg.GetEndpoint()}))
	}

	if cfg.GetClientId() != "" {
		o = append(o, mqtt.WithClientId(cfg.GetClientId()))
	}

	if cfg.GetUsername() != "" && cfg.GetPassword() != "" {
		o = append(o, mqtt.WithAuth(cfg.GetUsername(), cfg.GetPassword()))
	}

	o = append(o, mqtt.WithCleanSession(cfg.GetCleanSession()))

	var brokerOpts []broker.Option
	if cfg.GetKeepAlive().AsDuration() > 0 {
		brokerOpts = append(brokerOpts, mqttBroker.WithKeepAlive(cfg.GetKeepAlive().AsDuration()))
	}
	if cfg.GetAutoReconnect() {
		brokerOpts = append(brokerOpts, mqttBroker.WithAutoReconnect(true))
		brokerOpts = append(brokerOpts, mqttBroker.WithResumeSubs(cfg.GetResumeSubs()))
		if cfg.GetMaxReconnectInterval().AsDuration() > 0 {
			brokerOpts = append(brokerOpts, mqttBroker.WithMaxReconnectInterval(cfg.GetMaxReconnectInterval().AsDuration()))
		}
		if cfg.GetConnectRetryInterval().AsDuration() > 0 {
			brokerOpts = append(brokerOpts, mqttBroker.WithConnectRetryInterval(cfg.GetConnectRetryInterval().AsDuration()))
		}
	}
	if cfg.GetOrderMatters() {
		brokerOpts = append(brokerOpts, mqttBroker.WithOrderMatters(true))
	}
	if cfg.GetWriteTimeout().AsDuration() > 0 {
		brokerOpts = append(brokerOpts, mqttBroker.WithWriteTimeout(cfg.GetWriteTimeout().AsDuration()))
	}
	if len(brokerOpts) > 0 {
		o = append(o, mqtt.WithBrokerOptions(brokerOpts...))
	}

	if cfg.Tls != nil {
		tlsCfg, err := loadClientTlsConfig(cfg.Tls)
		if err != nil {
			log.Errorf("mqtt: load tls config failed: %v", err)
		} else {
			o = append(o, mqtt.WithTLSConfig(tlsCfg))
		}
	}

	if opts != nil {
		o = append(o, opts...)
	}

	srv := mqtt.NewServer(o...)

	return srv
}

func loadClientTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var out *tls.Config
	var err error

	if cfg.File != nil {
		if out, err = tlsUtils.LoadClientTlsConfigFile(
			cfg.File.GetKeyPath(),
			cfg.File.GetCertPath(),
			cfg.File.GetCaPath(),
		); err != nil {
			return nil, err
		}
	} else if cfg.Config != nil {
		if out, err = tlsUtils.LoadClientTlsConfigString(
			cfg.Config.GetKeyPem(),
			cfg.Config.GetCertPem(),
			cfg.Config.GetCaPem(),
		); err != nil {
			return nil, err
		}
	}

	if out != nil && cfg.GetInsecureSkipVerify() {
		out.InsecureSkipVerify = true
	}

	return out, nil
}
