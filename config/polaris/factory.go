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
func NewConfigSource(_ *conf.RemoteConfig) (config.Source, error) {
	configApi, err := polarisApi.NewConfigAPI()
	if err != nil {
		log.Fatal(err)
	}

	var opts []Option
	opts = append(opts, WithNamespace("default"))
	opts = append(opts, WithFileGroup("default"))
	opts = append(opts, WithFileName("default.yaml"))

	src, err := New(configApi, opts...)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return src, nil
}
