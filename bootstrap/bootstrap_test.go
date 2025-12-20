package bootstrap

import (
	"testing"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/stretchr/testify/assert"
	v1 "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func initApp(logger log.Logger, registrar registry.Registrar, _ *v1.Bootstrap) (*kratos.App, func(), error) {
	app := NewApp(logger, registrar)
	return app, func() {
	}, nil
}

func TestBootstrap(t *testing.T) {
	serviceName := "test"
	version := "v0.0.1"
	err := Bootstrap(initApp, &serviceName, &version)
	assert.Nil(t, err)
}

func TestCustomBootstrap(t *testing.T) {
	//customCfg := &CustomConfig{}
	//bConfig.RegisterConfig(customCfg)
	//
	//Bootstrap(func(logger log.Logger, registrar registry.Registrar, bootstrap *v1.Bootstrap) (*kratos.App, func(), error) {
	//	return initApp(logger, registrar, bootstrap, customCfg)
	//}, trans.Ptr("test"), trans.Ptr("1.0.0"))
}
