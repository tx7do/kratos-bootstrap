package nacos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewNacosRegistry(t *testing.T) {
	cfg := conf.Registry{
		Nacos: &conf.Registry_Nacos{
			Address: "127.0.0.1",
			Port:    8848,
		},
	}

	reg := NewRegistry(&cfg)
	assert.NotNil(t, reg)
}
