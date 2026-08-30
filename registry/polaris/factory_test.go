package polaris

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewPolarisRegistry(t *testing.T) {
	cfg := conf.Registry{
		Polaris: &conf.Registry_Polaris{
			Address:       "127.0.0.1",
			Port:          8091,
			InstanceCount: 0,
			Namespace:     "default",
			Service:       "DiscoverEchoServer",
			Token:         "",
			Weight:        100,
			Priority:      1,
			Healthy:       protoBool(true),
			Isolate:       protoBool(false),
			Heartbeat:     protoBool(true),
			Ttl:           5,
			Protocol:      "grpc",
			Timeout:       durationpb.New(3 * time.Second),
			RetryCount:    2,
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)

	assert.Equal(t, "default", reg.opt.Namespace)
	assert.Equal(t, 100, reg.opt.Weight)
	assert.Equal(t, 1, reg.opt.Priority)
	assert.Equal(t, 5, reg.opt.TTL)
	assert.Equal(t, 3*time.Second, reg.opt.Timeout)
	assert.Equal(t, 2, reg.opt.RetryCount)
}

func protoBool(v bool) *bool { return &v }
