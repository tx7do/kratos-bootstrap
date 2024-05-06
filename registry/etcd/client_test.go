package etcd

import (
	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"testing"
)

func TestNewEtcdRegistry(t *testing.T) {
	var cfg conf.Registry
	cfg.Etcd.Endpoints = []string{"127.0.0.1:2379"}

	reg := NewRegistry(&cfg)
	assert.Nil(t, reg)
}
