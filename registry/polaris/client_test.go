package polaris

import (
	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"testing"
)

func TestNewPolarisRegistry(t *testing.T) {
	var cfg conf.Registry
	cfg.Polaris.Address = "127.0.0.1"
	cfg.Polaris.Port = 8091
	cfg.Polaris.InstanceCount = 5
	cfg.Polaris.Namespace = "default"
	cfg.Polaris.Service = "DiscoverEchoServer"
	cfg.Polaris.Token = ""

	reg := NewPolarisRegistry(&cfg)
	assert.Nil(t, reg)
}
