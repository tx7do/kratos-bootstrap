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
		TimeoutMs:       c.Nacos.TimeoutMs,            // http请求超时时间，单位毫秒
		BeatInterval:    c.Nacos.BeatInterval,         // 心跳间隔时间，单位毫秒
		UpdateThreadNum: int(c.Nacos.UpdateThreadNum), // 更新服务的线程数

		NotLoadCacheAtStart:  c.Nacos.NotLoadCacheAtStart,  // 在启动时不读取本地缓存数据，true--不读取，false--读取
		UpdateCacheWhenEmpty: c.Nacos.UpdateCacheWhenEmpty, // 当服务列表为空时是否更新本地缓存，true--更新,false--不更新

		LogLevel: c.Nacos.LogLevel, // 日志级别
		CacheDir: c.Nacos.CacheDir, // 缓存目录
		LogDir:   c.Nacos.LogDir,   // 日志目录

		Username: c.Nacos.Username, // 用户名
		Password: c.Nacos.Password, // 密码

		NamespaceId: c.Nacos.NamespaceId, // 命名空间ID
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
