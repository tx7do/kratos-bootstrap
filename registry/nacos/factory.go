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

	in := c.Nacos

	serverConf := *nacosConstant.NewServerConfig(in.GetAddress(), in.GetPort())
	if in.GetScheme() != "" {
		serverConf.Scheme = in.GetScheme()
	}
	if in.GetContextPath() != "" {
		serverConf.ContextPath = in.GetContextPath()
	}
	if in.GetGrpcPort() > 0 {
		serverConf.GrpcPort = in.GetGrpcPort()
	}
	srvConf := []nacosConstant.ServerConfig{serverConf}

	cliConf := nacosConstant.ClientConfig{
		NamespaceId: in.GetNamespaceId(),
		RegionId:    in.GetRegionId(), // 地域ID
		AppName:     in.GetAppName(),
		AppKey:      in.GetAppKey(),

		TimeoutMs:    uint64(in.GetTimeout().AsDuration().Milliseconds()), // http请求超时时间，单位毫秒
		BeatInterval: in.GetBeatInterval().AsDuration().Milliseconds(),    // 心跳间隔时间，单位毫秒

		ListenInterval: uint64(in.GetListenInterval().AsDuration().Milliseconds()), // 监听间隔时间（已废弃），单位毫秒

		UpdateThreadNum:      int(in.GetUpdateThreadNum()), // 更新服务的线程数
		LogLevel:             in.GetLogLevel(),
		CacheDir:             in.GetCacheDir(),             // 缓存目录
		LogDir:               in.GetLogDir(),               // 日志目录
		NotLoadCacheAtStart:  in.GetNotLoadCacheAtStart(),  // 在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: in.GetUpdateCacheWhenEmpty(), // 当服务列表为空时是否更新本地缓存，true--更新,false--不更新

		Username: in.GetUsername(),
		Password: in.GetPassword(),

		OpenKMS: in.GetOpenKms(), // 是否开启KMS加密

		AccessKey: in.GetAccessKey(), // 阿里云AccessKey
		SecretKey: in.GetSecretKey(), // 阿里云SecretKey

		ContextPath: in.GetContextPath(),

		Endpoint:           in.GetEndpoint(),           // 地址服务器endpoint
		DisableUseSnapShot: in.GetDisableUseSnapShot(), // 请求服务端失败时禁用本地快照兜底
		AppendToStdout:     in.GetAppendToStdout(),     // 日志追加输出到标准输出
		AsyncUpdateService: in.GetAsyncUpdateService(), // 异步订阅更新服务列表
	}

	cli, err := nacosClients.NewNamingClient(
		nacosVo.NacosClientParam{
			ClientConfig:  &cliConf,
			ServerConfigs: srvConf,
		},
	)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	var opts []Option
	if in.GetGroup() != "" {
		opts = append(opts, WithGroup(in.GetGroup()))
	}
	if in.GetCluster() != "" {
		opts = append(opts, WithCluster(in.GetCluster()))
	}
	if in.GetWeight() > 0 {
		opts = append(opts, WithWeight(in.GetWeight()))
	}
	if in.GetKind() != "" {
		opts = append(opts, WithDefaultKind(in.GetKind()))
	}

	reg := New(cli, opts...)

	return reg, nil
}

func NewDiscovery(c *conf.Registry) (registry.Discovery, error) {
	return NewRegistry(c)
}

func NewRegistrar(c *conf.Registry) (registry.Registrar, error) {
	return NewRegistry(c)
}
