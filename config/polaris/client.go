package polaris

import (
	polarisKratos "github.com/go-kratos/kratos/contrib/config/polaris/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	polarisApi "github.com/polarismesh/polaris-go"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewConfigSource 创建一个远程配置源 - Polaris
func NewConfigSource(_ *conf.RemoteConfig) config.Source {
	configApi, err := polarisApi.NewConfigAPI()
	if err != nil {
		log.Fatal(err)
	}

	var opts []polarisKratos.Option
	opts = append(opts, polarisKratos.WithNamespace("default"))
	opts = append(opts, polarisKratos.WithFileGroup("default"))
	opts = append(opts, polarisKratos.WithFileName("default.yaml"))

	source, err := polarisKratos.New(configApi, opts...)
	if err != nil {
		log.Fatal(err)
	}

	return source
	return nil
}
