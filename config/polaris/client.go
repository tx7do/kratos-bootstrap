package polaris

import (
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	"github.com/go-kratos/kratos/v2/config"
)

// NewConfigSource 创建一个远程配置源 - Polaris
func NewConfigSource(_ *conf.RemoteConfig) config.Source {
	//configApi, err := polarisApi.NewConfigAPI()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//var opts []polarisKratos.Option
	//opts = append(opts, polarisKratos.WithNamespace("default"))
	//opts = append(opts, polarisKratos.WithFileGroup("default"))
	//opts = append(opts, polarisKratos.WithFileName("default.yaml"))
	//
	//source, err := polarisKratos.New(configApi, opts...)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//return source
	return nil
}
