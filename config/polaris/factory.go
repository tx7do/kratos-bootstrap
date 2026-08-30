package polaris

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"

	polaris "github.com/polarismesh/polaris-go"
	polarisConfig "github.com/polarismesh/polaris-go/pkg/config"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	bConfig "github.com/tx7do/kratos-bootstrap/config"
)

func init() {
	bConfig.MustRegisterFactory(bConfig.TypePolaris, NewConfigSource)
}

// NewConfigSource 创建一个远程配置源 - Polaris
func NewConfigSource(cfg *conf.RemoteConfig) (config.Source, error) {
	if cfg == nil || cfg.Polaris == nil {
		return nil, nil
	}

	in := cfg.Polaris

	var configApi polaris.ConfigAPI
	var err error
	switch {
	case in.GetConfigFile() != "":
		if configApi, err = polaris.NewConfigAPIByFile(in.GetConfigFile()); err != nil {
			log.Errorf("fail to create configAPI by file, err is %v", err)
			return nil, err
		}
	case in.GetAddress() != "":
		address := in.GetAddress()
		if in.GetPort() > 0 {
			address = fmt.Sprintf("%s:%d", in.GetAddress(), in.GetPort())
		}
		configuration := polarisConfig.NewDefaultConfiguration([]string{address})
		if configApi, err = polaris.NewConfigAPIByConfig(configuration); err != nil {
			log.Errorf("fail to create configAPI by address, err is %v", err)
			return nil, err
		}
	default:
		if configApi, err = polaris.NewConfigAPI(); err != nil {
			log.Errorf("%v", err)
			return nil, err
		}
	}

	var opts []Option
	opts = append(opts, WithNamespace(in.GetNamespace()))
	opts = append(opts, WithFileGroup(in.GetFileGroup()))
	opts = append(opts, WithFileName(in.GetFileName()))

	src, err := New(configApi, opts...)
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}

	return src, nil
}
