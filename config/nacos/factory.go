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

	srvConf := []nacosConstant.ServerConfig{
		*nacosConstant.NewServerConfig(c.Nacos.Address, c.Nacos.Port),
	}

	cliConf := nacosConstant.ClientConfig{
		TimeoutMs:            10 * 1000, // http请求超时时间，单位毫秒
		BeatInterval:         5 * 1000,  // 心跳间隔时间，单位毫秒
		UpdateThreadNum:      20,        // 更新服务的线程数
		LogLevel:             "debug",
		CacheDir:             "../../configs/cache", // 缓存目录
		LogDir:               "../../configs/log",   // 日志目录
		NotLoadCacheAtStart:  true,                  // 在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: true,                  // 当服务列表为空时是否更新本地缓存，true--更新,false--不更新
		Username:             c.Nacos.Username,      //用户名
		Password:             c.Nacos.Password,      //密码
		NamespaceId:          c.Nacos.NamespaceId,   //命名空间ID
	}

	nacosClient, err := nacosClients.NewConfigClient(
		nacosVo.NacosClientParam{
			ClientConfig:  &cliConf,
			ServerConfigs: srvConf,
		},
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return New(nacosClient,
		WithGroup(getConfigKey(c.Nacos.Key, false)),
		WithDataID("bootstrap.yaml"),
	), nil
}
