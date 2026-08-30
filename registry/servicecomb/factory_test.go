package servicecomb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewServicecombRegistry(t *testing.T) {
	cfg := conf.Registry{
		Servicecomb: &conf.Registry_Servicecomb{
			Endpoints: []string{"127.0.0.1:30100"},
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}

func TestNewServicecombRegistryWithFullParams(t *testing.T) {
	cfg := conf.Registry{
		Servicecomb: &conf.Registry_Servicecomb{
			Endpoints:       []string{"127.0.0.1:30100"},
			Timeout:         durationpb.New(15 * time.Second),
			EnableAuth:      true,
			Username:        "user",
			Password:        "password",
			TokenExpiration: durationpb.New(time.Hour),
			Verbose:         true,
			Compressed:      true,
		},
	}

	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}
