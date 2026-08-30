package consul

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewConsulRegistry(t *testing.T) {
	cfg := conf.Registry{
		Consul: &conf.Registry_Consul{
			Scheme:      "http",
			Address:     "localhost:8500",
			HealthCheck: false,
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}

func TestNewConsulRegistryWithFullParams(t *testing.T) {
	cfg := conf.Registry{
		Consul: &conf.Registry_Consul{
			Scheme:                         "http",
			Address:                        "localhost:8500",
			HealthCheck:                    true,
			Datacenter:                     "dc1",
			Token:                          "token",
			Namespace:                      "ns",
			Partition:                      "default",
			PathPrefix:                     "/consul",
			BasicAuth:                      &conf.Registry_Consul_BasicAuth{Username: "u", Password: "p"},
			WaitTime:                       durationpb.New(10 * time.Second),
			Timeout:                        durationpb.New(5 * time.Second),
			Heartbeat:                      protoBool(false),
			HealthCheckInterval:            protoInt32(5),
			DeregisterCriticalServiceAfter: protoInt32(60),
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)

	r := reg
	assert.Equal(t, 5*time.Second, r.timeout)
	assert.Equal(t, 5, r.cli.healthcheckInterval)
	assert.Equal(t, 60, r.cli.deregisterCriticalServiceAfter)
	assert.False(t, r.cli.heartbeat)
}

func protoBool(v bool) *bool    { return &v }
func protoInt32(v int32) *int32 { return &v }
