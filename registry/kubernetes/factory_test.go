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
	assert.Nil(t, reg)

	cfg.Kubernetes = &conf.Registry_Kubernetes{
		Kubeconfig: "/nonexistent/kubeconfig",
	}
	// kubeconfig不存在时应返回错误
	_, err = NewRegistry(&cfg)
	assert.NotNil(t, err)
}
