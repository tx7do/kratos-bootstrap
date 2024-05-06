package kubernetes

import (
	"github.com/stretchr/testify/assert"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"testing"
)

func TestNewKubernetesRegistry(t *testing.T) {
	var cfg conf.Registry
	reg := NewRegistry(&cfg)
	assert.Nil(t, reg)
}
