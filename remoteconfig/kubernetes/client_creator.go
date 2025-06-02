package kubernetes

import (
	"path/filepath"

	k8sUtil "k8s.io/client-go/util/homedir"

	"github.com/go-kratos/kratos/v2/config"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewConfigSource 创建一个远程配置源 - Kubernetes
func NewConfigSource(c *conf.RemoteConfig) config.Source {
	if c == nil || c.Kubernetes == nil {
		return nil
	}

	source := NewSource(
		Namespace(c.Kubernetes.Namespace),
		LabelSelector(""),
		KubeConfig(filepath.Join(k8sUtil.HomeDir(), ".kube", "config")),
	)
	return source
}
