package zookeeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewZooKeeperRegistry(t *testing.T) {
	cfg := conf.Registry{
		Zookeeper: &conf.Registry_ZooKeeper{
			Endpoints: []string{"127.0.0.1:2181"},
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}
