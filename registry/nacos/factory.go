package nacos

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	nacosClients "github.com/nacos-group/nacos-sdk-go/v2/clients"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nacosVo "github.com/nacos-group/nacos-sdk-go/v2/vo"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	baseRegistry "github.com/tx7do/kratos-bootstrap/registry"
)

func init() {
	_ = baseRegistry.RegisterDiscoveryFactory(baseRegistry.Nacos, NewDiscovery)
	_ = baseRegistry.RegisterRegistrarFactory(baseRegistry.Nacos, NewRegistrar)
}

// NewRegistry 创建一个注册发现客户端 - Nacos
func NewRegistry(c *conf.Registry) (*Registry, error) {
	if c == nil || c.Nacos == nil {
		return nil, nil
	}

	srvConf := []nacosConstant.ServerConfig{
		*nacosConstant.NewServerConfig(c.Nacos.Address, c.Nacos.Port),
	}

	cliConf := nacosConstant.ClientConfig{
		NamespaceId: c.Nacos.NamespaceId,
		RegionId:    c.Nacos.RegionId, // 地域ID
		AppName:     c.Nacos.AppName,
		AppKey:      c.Nacos.AppKey,

		TimeoutMs:    uint64(c.Nacos.Timeout.AsDuration().Milliseconds()), // http请求超时时间，单位毫秒
		BeatInterval: c.Nacos.BeatInterval.AsDuration().Milliseconds(),    // 心跳间隔时间，单位毫秒

		UpdateThreadNum:      int(c.Nacos.UpdateThreadNum), // 更新服务的线程数
		LogLevel:             c.Nacos.LogLevel,
		CacheDir:             c.Nacos.CacheDir,             // 缓存目录
		LogDir:               c.Nacos.LogDir,               // 日志目录
		NotLoadCacheAtStart:  c.Nacos.NotLoadCacheAtStart,  // 在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: c.Nacos.UpdateCacheWhenEmpty, // 当服务列表为空时是否更新本地缓存，true--更新,false--不更新

		Username: c.Nacos.Username,
		Password: c.Nacos.Password,

		OpenKMS: c.Nacos.OpenKms, // 是否开启KMS加密

		AccessKey: c.Nacos.AccessKey, // 阿里云AccessKey
		SecretKey: c.Nacos.SecretKey, // 阿里云SecretKey

		ContextPath: c.Nacos.ContextPath,
	}

	cli, err := nacosClients.NewNamingClient(
		nacosVo.NacosClientParam{
			ClientConfig:  &cliConf,
			ServerConfigs: srvConf,
		},
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	reg := New(cli)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
