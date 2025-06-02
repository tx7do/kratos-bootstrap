package etcd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewEtcdRegistry(t *testing.T) {
	cfg := conf.Registry{
		Etcd: &conf.Registry_Etcd{
			Endpoints: []string{"127.0.0.1:2379"},
		},
	}

	reg := NewRegistry(&cfg)
	assert.NotNil(t, reg)
}
