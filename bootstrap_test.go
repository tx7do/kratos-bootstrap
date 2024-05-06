package bootstrap

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	v1 "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"testing"
)

func initApp(logger log.Logger, registrar registry.Registrar, bootstrap *v1.Bootstrap) (*kratos.App, func(), error) {
	app := NewApp(logger, registrar)
	return app, func() {
	}, nil
}

func TestBootstrap(t *testing.T) {
	Bootstrap(initApp)
}
