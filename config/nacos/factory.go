package nacos

import (
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	nacosClients "github.com/nacos-group/nacos-sdk-go/v2/clients"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nacosVo "github.com/nacos-group/nacos-sdk-go/v2/vo"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

func init() {
	bConfig.MustRegisterFactory(bConfig.TypeNacos, NewConfigSource)
}

// NewConfigSource 创建一个远程配置源 - Nacos
func NewConfigSource(c *conf.RemoteConfig) (config.Source, error) {
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
		TimeoutMs:       in.GetTimeoutMs(),            // http请求超时时间，单位毫秒
		BeatInterval:    in.GetBeatInterval(),         // 心跳间隔时间，单位毫秒
		UpdateThreadNum: int(in.GetUpdateThreadNum()), // 更新服务的线程数

		NotLoadCacheAtStart:  in.GetNotLoadCacheAtStart(),  // 在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: in.GetUpdateCacheWhenEmpty(), // 当配置为空时更新缓存

		LogLevel: in.GetLogLevel(), // 日志级别
		CacheDir: in.GetCacheDir(), // 缓存目录
		LogDir:   in.GetLogDir(),   // 日志目录

		Username: in.GetUsername(), // 用户名
		Password: in.GetPassword(), // 密码

		NamespaceId: in.GetNamespaceId(), // 命名空间ID

		Endpoint:           in.GetEndpoint(),           // 地址服务器endpoint
		ContextPath:        in.GetContextPath(),        // 上下文路径
		AppendToStdout:     in.GetAppendToStdout(),     // 日志追加输出到标准输出
		DisableUseSnapShot: in.GetDisableUseSnapShot(), // 请求服务端失败时禁用本地快照兜底
	}

	if tlsConf := in.GetTls(); tlsConf != nil {
		cliConf.TLSCfg = nacosConstant.TLSConfig{
			Appointed: true,
			Enable:    true,
			TrustAll:  tlsConf.GetInsecureSkipVerify(),
		}
		if file := tlsConf.GetFile(); file != nil {
			cliConf.TLSCfg.CaFile = file.GetCaPath()
			cliConf.TLSCfg.CertFile = file.GetCertPath()
			cliConf.TLSCfg.KeyFile = file.GetKeyPath()
		}
	}

	nacosClient, err := nacosClients.NewConfigClient(
		nacosVo.NacosClientParam{
			ClientConfig:  &cliConf,
			ServerConfigs: srvConf,
		},
	)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	var group string
	if c.Nacos.GetGroup() != "" {
		group = c.Nacos.GetGroup()
	} else {
		group = DefaultGroup
	}

	var dataID string
	if c.Nacos.GetDataId() != "" {
		dataID = c.Nacos.GetDataId()
	} else {
		dataID = DefaultDataID
	}

	return New(nacosClient,
		WithGroup(group),
		WithDataID(dataID),
	), nil
}
