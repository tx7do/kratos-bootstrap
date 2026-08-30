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

	in := c.Kubernetes

	opts := []Option{
		WithNamespace(in.GetNamespace()),
		WithLabelSelector(in.GetLabelSelector()),
		WithFieldSelector(in.GetFieldSelector()),
		WithMaster(in.GetMaster()),
		WithKubeConfig(in.GetKubeConfig()),
	}
	if in.Qps != nil && in.GetQps() > 0 {
		opts = append(opts, WithQPS(int(in.GetQps())))
	}
	if in.Burst != nil && in.GetBurst() > 0 {
		opts = append(opts, WithBurst(int(in.GetBurst())))
	}
	if in.GetTimeout().AsDuration() > 0 {
		opts = append(opts, WithRequestTimeout(in.GetTimeout().AsDuration()))
	}

	return NewSource(opts...), nil
}
