package eureka

import (
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewRegistry 创建一个注册发现客户端 - Eureka
func NewRegistry(c *conf.Registry) *Registry {
	if c == nil || c.Eureka == nil {
		return nil
	}

	var opts []Option
	if c.Eureka.HeartbeatInterval != nil {
		opts = append(opts, WithHeartbeat(c.Eureka.HeartbeatInterval.AsDuration()))
	}
	if c.Eureka.RefreshInterval != nil {
		opts = append(opts, WithRefresh(c.Eureka.RefreshInterval.AsDuration()))
	}
	if c.Eureka.Path != "" {
		opts = append(opts, WithEurekaPath(c.Eureka.Path))
	}

	var err error
	var reg *Registry
	if reg, err = New(c.Eureka.Endpoints, opts...); err != nil {
		log.Fatal(err)
	}

	return reg
}
