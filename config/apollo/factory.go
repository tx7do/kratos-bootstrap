package apollo

import (
	"github.com/go-kratos/kratos/v2/config"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

func init() {
	bConfig.MustRegisterFactory(bConfig.TypeApollo, NewConfigSource)
}

// NewConfigSource 创建一个远程配置源 - Apollo
func NewConfigSource(cfg *conf.RemoteConfig) (config.Source, error) {
	if cfg == nil || cfg.Apollo == nil {
		return nil, nil
	}

	source := NewSource(
		WithAppID(cfg.Apollo.AppId),
		WithCluster(cfg.Apollo.Cluster),
		WithEndpoint(cfg.Apollo.Endpoint),
		WithNamespace(cfg.Apollo.Namespace),
		WithSecret(cfg.Apollo.Secret),
		WithEnableBackup(),
	)
	return source, nil
}
