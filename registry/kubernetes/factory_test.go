package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func TestNewKubernetesRegistry(t *testing.T) {
	var cfg conf.Registry
	reg, err := NewRegistry(&cfg)
	assert.Nil(t, err)
	assert.NotNil(t, reg)
}
