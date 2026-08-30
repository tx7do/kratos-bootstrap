package eureka

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewEurekaRegistry(t *testing.T) {
	cfg := conf.Registry{
		Eureka: &conf.Registry_Eureka{
			Endpoints: []string{"https://127.0.0.1:18761"},
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}

func TestNewEurekaRegistryWithFullParams(t *testing.T) {
	cfg := conf.Registry{
		Eureka: &conf.Registry_Eureka{
			Endpoints:         []string{"https://127.0.0.1:18761"},
			HeartbeatInterval: durationpb.New(15 * time.Second),
			RefreshInterval:   durationpb.New(20 * time.Second),
			Path:              "eureka/v2",
			MaxRetry:          5,
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)

	assert.Equal(t, 15*time.Second, reg.heartbeatInterval)
	assert.Equal(t, 20*time.Second, reg.refreshInterval)
	assert.Equal(t, "eureka/v2", reg.eurekaPath)
	assert.Equal(t, 5, reg.maxRetry)
}
