package eureka

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
