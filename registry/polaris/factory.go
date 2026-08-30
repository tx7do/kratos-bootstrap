package polaris

import (
	"fmt"

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

	in := c.Polaris

	var err error

	var consumer polarisApi.ConsumerAPI
	switch {
	case in.GetConfigFile() != "":
		if consumer, err = polarisApi.NewConsumerAPIByFile(in.GetConfigFile()); err != nil {
			log.Errorf("fail to create consumerAPI by file, err is %v", err)
			return nil, err
		}
	case in.GetAddress() != "":
		address := in.GetAddress()
		if in.GetPort() > 0 {
			address = fmt.Sprintf("%s:%d", in.GetAddress(), in.GetPort())
		}
		if consumer, err = polarisApi.NewConsumerAPIByAddress(address); err != nil {
			log.Errorf("fail to create consumerAPI by address, err is %v", err)
			return nil, err
		}
	default:
		if consumer, err = polarisApi.NewConsumerAPI(); err != nil {
			log.Errorf("fail to create consumerAPI, err is %v", err)
			return nil, err
		}
	}

	var provider polarisApi.ProviderAPI
	provider = polarisApi.NewProviderAPIByContext(consumer.SDKContext())

	var opts []Option
	if in.GetNamespace() != "" {
		opts = append(opts, WithNamespace(in.GetNamespace()))
	}
	if in.GetToken() != "" {
		opts = append(opts, WithServiceToken(in.GetToken()))
	}
	if in.GetWeight() > 0 {
		opts = append(opts, WithWeight(int(in.GetWeight())))
	}
	if in.GetPriority() > 0 {
		opts = append(opts, WithPriority(int(in.GetPriority())))
	}
	if in.GetProtocol() != "" {
		opts = append(opts, WithProtocol(in.GetProtocol()))
	}
	if in.Healthy != nil {
		opts = append(opts, WithHealthy(in.GetHealthy()))
	}
	if in.Isolate != nil {
		opts = append(opts, WithIsolate(in.GetIsolate()))
	}
	if in.Heartbeat != nil {
		opts = append(opts, WithHeartbeat(in.GetHeartbeat()))
	}
	if in.GetTimeout().AsDuration() > 0 {
		opts = append(opts, WithTimeout(in.GetTimeout().AsDuration()))
	}
	if in.GetRetryCount() > 0 {
		opts = append(opts, WithRetryCount(int(in.GetRetryCount())))
	}
	if in.GetTtl() > 0 {
		opts = append(opts, WithTTL(int(in.GetTtl())))
	} else if in.Heartbeat == nil || in.GetHeartbeat() {
		// 心跳默认开启，TTL为0会导致心跳ticker panic，兜底一个默认值
		opts = append(opts, WithTTL(5))
	}

	log.Infof("start to register instances, count %d", in.GetInstanceCount())

	var resp *polarisModel.InstanceRegisterResponse
	for i := 0; i < (int)(in.GetInstanceCount()); i++ {
		registerRequest := &polarisApi.InstanceRegisterRequest{}
		registerRequest.Service = in.GetService()
		registerRequest.Namespace = in.GetNamespace()
		registerRequest.Host = in.GetAddress()
		registerRequest.Port = (int)(in.GetPort()) + i
		registerRequest.ServiceToken = in.GetToken()
		registerRequest.SetHealthy(true)
		if resp, err = provider.RegisterInstance(registerRequest); err != nil {
			log.Errorf("fail to register instance %d, err is %v", i, err)
			return nil, err
		}

		log.Infof("register instance %d response: instanceId %s", i, resp.InstanceID)
	}

	reg := New(provider, consumer, opts...)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
