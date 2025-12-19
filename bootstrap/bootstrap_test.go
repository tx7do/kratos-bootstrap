package bootstrap

import (
	"testing"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

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
	Bootstrap(initApp, &serviceName, &version)
}

type CustomConfig struct {
	Cfg string `protobuf:"bytes,1,opt,name=cfg,proto3" json:"cfg,omitempty"`
}

func initAppEx(logger log.Logger, registrar registry.Registrar, _ *v1.Bootstrap, _ *CustomConfig) (*kratos.App, func(), error) {
	app := NewApp(logger, registrar)
	return app, func() {
	}, nil
}

func TestCustomBootstrap(t *testing.T) {
	customCfg := &CustomConfig{}
	RegisterConfig(customCfg)

	AppInfo.Name = "test"
	AppInfo.Version = "v0.0.1"

	// bootstrap
	cfg, ll, reg := DoBootstrap(AppInfo)

	// init app
	app, cleanup, err := initAppEx(ll, reg, cfg, customCfg)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// run the app.
	if err = app.Run(); err != nil {
		panic(err)
	}
}
