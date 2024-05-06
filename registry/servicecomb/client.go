package servicecomb

import (
	"github.com/go-kratos/kratos/v2/log"

	servicecombClient "github.com/go-chassis/sc-client"
	servicecombKratos "github.com/go-kratos/kratos/contrib/registry/servicecomb/v2"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewRegistry 创建一个注册发现客户端 - Servicecomb
func NewRegistry(c *conf.Registry) *servicecombKratos.Registry {
	cfg := servicecombClient.Options{
		Endpoints: c.Servicecomb.Endpoints,
	}

	var cli *servicecombClient.Client
	var err error
	if cli, err = servicecombClient.NewClient(cfg); err != nil {
		log.Fatal(err)
	}

	reg := servicecombKratos.NewRegistry(cli)

	return reg
}
