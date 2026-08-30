package zookeeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewZooKeeperRegistry(t *testing.T) {
	cfg := conf.Registry{
		Zookeeper: &conf.Registry_ZooKeeper{
			Endpoints: []string{"127.0.0.1:2181"},
			Timeout:   durationpb.New(5 * time.Second),
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}

func TestNewZooKeeperRegistryWithFullParams(t *testing.T) {
	cfg := conf.Registry{
		Zookeeper: &conf.Registry_ZooKeeper{
			Endpoints: []string{"127.0.0.1:2181"},
			Timeout:   durationpb.New(5 * time.Second),
			Namespace: "/my-services",
			Username:  "user",
			Password:  "password",
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)

	assert.Equal(t, "/my-services", reg.opts.namespace)
	assert.Equal(t, "user", reg.opts.user)
	assert.Equal(t, "password", reg.opts.password)
}
