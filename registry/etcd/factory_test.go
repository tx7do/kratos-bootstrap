package etcd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewEtcdRegistry(t *testing.T) {
	cfg := conf.Registry{
		Etcd: &conf.Registry_Etcd{
			Endpoints: []string{"127.0.0.1:2379"},
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}

func TestNewEtcdRegistryWithFullParams(t *testing.T) {
	cfg := conf.Registry{
		Etcd: &conf.Registry_Etcd{
			Endpoints:            []string{"127.0.0.1:2379"},
			Username:             "root",
			Password:             "password",
			DialTimeout:          durationpb.New(5 * time.Second),
			AutoSyncInterval:     durationpb.New(time.Minute),
			DialKeepAliveTime:    durationpb.New(30 * time.Second),
			DialKeepAliveTimeout: durationpb.New(10 * time.Second),
			PermitWithoutStream:  true,
			Namespace:            "/my-services",
			RegisterTtl:          durationpb.New(30 * time.Second),
			MaxRetry:             3,
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)

	assert.Equal(t, "/my-services", reg.opts.namespace)
	assert.Equal(t, 30*time.Second, reg.opts.ttl)
	assert.Equal(t, 3, reg.opts.maxRetry)
}
