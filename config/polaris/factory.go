package polaris

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	polarisApi "github.com/polarismesh/polaris-go"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

func init() {
	bConfig.MustRegisterFactory(bConfig.TypeNacos, NewConfigSource)
}

// NewConfigSource 创建一个远程配置源 - Polaris
func NewConfigSource(cfg *conf.RemoteConfig) (config.Source, error) {
	if cfg == nil || cfg.Polaris == nil {
		return nil, nil
	}

	configApi, err := polarisApi.NewConfigAPI()
	if err != nil {
		log.Fatal(err)
	}

	var opts []Option
	opts = append(opts, WithNamespace(cfg.Polaris.Namespace))
	opts = append(opts, WithFileGroup(cfg.Polaris.FileGroup))
	opts = append(opts, WithFileName(cfg.Polaris.FileName))

	src, err := New(configApi, opts...)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return src, nil
}
