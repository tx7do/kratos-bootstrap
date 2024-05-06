package eureka

import (
	"github.com/go-kratos/kratos/v2/log"

	eurekaKratos "github.com/go-kratos/kratos/contrib/registry/eureka/v2"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewRegistry 创建一个注册发现客户端 - Eureka
func NewRegistry(c *conf.Registry) *eurekaKratos.Registry {
	var opts []eurekaKratos.Option
	opts = append(opts, eurekaKratos.WithHeartbeat(c.Eureka.HeartbeatInterval.AsDuration()))
	opts = append(opts, eurekaKratos.WithRefresh(c.Eureka.RefreshInterval.AsDuration()))
	opts = append(opts, eurekaKratos.WithEurekaPath(c.Eureka.Path))

	var err error
	var reg *eurekaKratos.Registry
	if reg, err = eurekaKratos.New(c.Eureka.Endpoints, opts...); err != nil {
		log.Fatal(err)
	}

	return reg
}
