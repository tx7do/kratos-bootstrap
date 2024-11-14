package apollo

import (
	apolloKratos "github.com/go-kratos/kratos/contrib/config/apollo/v2"
	"github.com/go-kratos/kratos/v2/config"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewConfigSource 创建一个远程配置源 - Apollo
func NewConfigSource(cfg *conf.RemoteConfig) config.Source {
	if cfg == nil || cfg.Apollo == nil {
		return nil
	}

	source := apolloKratos.NewSource(
		apolloKratos.WithAppID(cfg.Apollo.AppId),
		apolloKratos.WithCluster(cfg.Apollo.Cluster),
		apolloKratos.WithEndpoint(cfg.Apollo.Endpoint),
		apolloKratos.WithNamespace(cfg.Apollo.Namespace),
		apolloKratos.WithSecret(cfg.Apollo.Secret),
		apolloKratos.WithEnableBackup(),
	)
	return source
}
