package mqtt

import (
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

	if opts != nil {
		o = append(o, opts...)
	}

	srv := mqtt.NewServer(o...)

	return srv
}
