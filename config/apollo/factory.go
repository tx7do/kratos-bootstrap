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

	in := cfg.Apollo

	opts := []Option{
		WithAppID(in.GetAppId()),
		WithCluster(in.GetCluster()),
		WithEndpoint(in.GetEndpoint()),
		WithNamespace(in.GetNamespace()),
		WithSecret(in.GetSecret()),
		WithEnableBackup(),
	}
	if in.GetBackupConfigPath() != "" {
		opts = append(opts, WithBackupPath(in.GetBackupConfigPath()))
	}
	if in.GetLabel() != "" {
		opts = append(opts, WithLabel(in.GetLabel()))
	}
	if in.GetSyncServerTimeout().AsDuration() > 0 {
		opts = append(opts, WithSyncServerTimeout(int(in.GetSyncServerTimeout().AsDuration().Seconds())))
	}
	opts = append(opts, WithMustStart(in.GetMustStart()))

	return NewSource(opts...), nil
}
