package bootstrap

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2"
	"github.com/stretchr/testify/assert"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

func initApp(ctx *Context) (*kratos.App, func(), error) {
	app := NewApp(ctx)
	return app, func() {
	}, nil
}

func TestBootstrapWithNameVersion(t *testing.T) {
	serviceName := "test"
	version := "v0.0.1"

	ctx := NewContext(context.Background(), &conf.AppInfo{
		Project: "",
		AppId:   serviceName,
		Version: version,
	})

	err := RunApp(ctx, initApp)
	assert.Nil(t, err)
}

func TestNewInstanceId(t *testing.T) {
	instanceId := NewInstanceId("gowind-test-service", "1.0.0", "127.0.0.1", "8000")
	t.Logf("InstanceId: %s", instanceId)
}
