package servicecomb

import (
	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"testing"
)

func TestNewServicecombRegistry(t *testing.T) {
	var cfg conf.Registry
	cfg.Servicecomb.Endpoints = []string{"127.0.0.1:30100"}

	reg := NewRegistry(&cfg)
	assert.Nil(t, reg)
}
