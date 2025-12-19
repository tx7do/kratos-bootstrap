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

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}
