package polaris

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	polarisApi "github.com/polarismesh/polaris-go/api"
	polarisModel "github.com/polarismesh/polaris-go/pkg/model"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.Polaris, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.Polaris, NewRegistrar)
}

// NewRegistry 创建一个注册发现客户端 - Polaris
func NewRegistry(c *conf.Registry) (*Registry, error) {
	if c == nil || c.Polaris == nil {
		return nil, nil
	}

	var err error

	var consumer polarisApi.ConsumerAPI
	if consumer, err = polarisApi.NewConsumerAPI(); err != nil {
		log.Fatalf("fail to create consumerAPI, err is %v", err)
		return nil, err
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
			return nil, err
		}

		log.Infof("register instance %d response: instanceId %s", i, resp.InstanceID)
	}

	reg := New(provider, consumer)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
