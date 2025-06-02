package polaris

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	polarisApi "github.com/polarismesh/polaris-go/api"
	polarisModel "github.com/polarismesh/polaris-go/pkg/model"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	r "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	r.RegisterRegistrarCreator(string(r.Polaris), func(c *conf.Registry) registry.Registrar {
		return NewRegistry(c)
	})
	r.RegisterDiscoveryCreator(string(r.Polaris), func(c *conf.Registry) registry.Discovery {
		return NewRegistry(c)
	})
}

// NewRegistry 创建一个注册发现客户端 - Polaris
func NewRegistry(c *conf.Registry) *Registry {
	if c == nil || c.Polaris == nil {
		return nil
	}

	var err error

	var consumer polarisApi.ConsumerAPI
	if consumer, err = polarisApi.NewConsumerAPI(); err != nil {
		log.Fatalf("fail to create consumerAPI, err is %v", err)
	}

	var provider polarisApi.ProviderAPI
	provider = polarisApi.NewProviderAPIByContext(consumer.SDKContext())

	log.Infof("start to register instances, count %d", c.Polaris.InstanceCount)

	var resp *polarisModel.InstanceRegisterResponse
	for i := 0; i < (int)(c.Polaris.InstanceCount); i++ {
		registerRequest := &polarisApi.InstanceRegisterRequest{}
		registerRequest.Service = c.Polaris.Service
		registerRequest.Namespace = c.Polaris.Namespace
		registerRequest.Host = c.Polaris.Address
		registerRequest.Port = (int)(c.Polaris.Port) + i
		registerRequest.ServiceToken = c.Polaris.Token
		registerRequest.SetHealthy(true)
		if resp, err = provider.RegisterInstance(registerRequest); err != nil {
			log.Fatalf("fail to register instance %d, err is %v", i, err)
		} else {
			log.Infof("register instance %d response: instanceId %s", i, resp.InstanceID)
		}
	}

	reg := New(provider, consumer)

	return reg
}
