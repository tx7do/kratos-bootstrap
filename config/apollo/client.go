package apollo

import (
	"github.com/go-kratos/kratos/v2/config"

	// apollo
	apolloKratos "github.com/go-kratos/kratos/contrib/config/apollo/v2"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewConfigSource 创建一个远程配置源 - Apollo
func NewConfigSource(c *conf.RemoteConfig) config.Source {
	source := apolloKratos.NewSource(
		apolloKratos.WithAppID(c.Apollo.AppId),
		apolloKratos.WithCluster(c.Apollo.Cluster),
		apolloKratos.WithEndpoint(c.Apollo.Endpoint),
		apolloKratos.WithNamespace(c.Apollo.Namespace),
		apolloKratos.WithSecret(c.Apollo.Secret),
		apolloKratos.WithEnableBackup(),
	)
	return source
}
