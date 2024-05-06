package nacos

import (
	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"testing"
)

func TestNewNacosRegistry(t *testing.T) {
	var cfg conf.Registry
	cfg.Nacos.Address = "127.0.0.1"
	cfg.Nacos.Port = 8848

	reg := NewRegistry(&cfg)
	assert.Nil(t, reg)
}
