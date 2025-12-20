package bootstrap

import (
	"testing"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/stretchr/testify/assert"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func initApp(logger log.Logger, registrar registry.Registrar, _ *conf.Bootstrap) (*kratos.App, func(), error) {
	app := NewApp(logger, registrar)
	return app, func() {
	}, nil
}

func TestBootstrapWithNameVersion(t *testing.T) {
	serviceName := "test"
	version := "v0.0.1"
	err := Bootstrap(func(ctx *Context) (app *kratos.App, cleanup func(), err error) {
		return initApp(ctx.Logger, ctx.Registrar, ctx.Config)
	}, &serviceName, &version)
	assert.Nil(t, err)
}

type CustomConfig struct {
}

func initAppEx(logger log.Logger, registrar registry.Registrar, _ *conf.Bootstrap, _ *CustomConfig) (*kratos.App, func(), error) {
	app := NewApp(logger, registrar)
	return app, func() {
	}, nil
}

func TestCustomBootstrap(t *testing.T) {
	//customCfg := &CustomConfig{}
	//bConfig.RegisterConfig(customCfg)
	//
	//Bootstrap(func(ctx *Context) (*kratos.App, func(), error) {
	//	return initAppEx(ctx.Logger, ctx.Registrar, ctx.Config, customCfg)
	//}, trans.Ptr("test"), trans.Ptr("1.0.0"))
}
