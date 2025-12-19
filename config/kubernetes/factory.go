package kubernetes

import (
	"github.com/go-kratos/kratos/v2/config"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

func init() {
	bConfig.MustRegisterFactory(bConfig.TypeKubernetes, NewConfigSource)
}

// NewConfigSource 创建一个远程配置源 - Kubernetes
func NewConfigSource(c *conf.RemoteConfig) (config.Source, error) {
	if c == nil || c.Kubernetes == nil {
		return nil, nil
	}

	source := NewSource(
		WithNamespace(c.Kubernetes.Namespace),
		WithLabelSelector(c.Kubernetes.LabelSelector),
		WithFieldSelector(c.Kubernetes.FieldSelector),
		WithMaster(c.Kubernetes.Master),
		WithKubeConfig(c.Kubernetes.KubeConfig),
	)
	return source, nil
}
