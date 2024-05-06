package consul

import (
	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"testing"
)

func TestNewConsulRegistry(t *testing.T) {
	var cfg conf.Registry
	cfg.Consul.Scheme = "http"
	cfg.Consul.Address = "localhost:8500"
	cfg.Consul.HealthCheck = false

	reg := NewRegistry(&cfg)
	assert.Nil(t, reg)
}
