package eureka

import (
	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"testing"
)

func TestNewEurekaRegistry(t *testing.T) {
	var cfg conf.Registry
	cfg.Eureka.Endpoints = []string{"https://127.0.0.1:18761"}

	reg := NewRegistry(&cfg)
	assert.Nil(t, reg)
}
