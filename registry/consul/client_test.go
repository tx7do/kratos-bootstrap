package consul

import (
	"testing"

	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewConsulRegistry(t *testing.T) {
	cfg := conf.Registry{
		Consul: &conf.Registry_Consul{
			Scheme:      "http",
			Address:     "localhost:8500",
			HealthCheck: false,
		},
	}

	reg := NewRegistry(&cfg)
	assert.NotNil(t, reg)
}
