package rpc

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"

	"google.golang.org/grpc"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	midRateLimit "github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"

	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/utils"
)

const defaultTimeout = 5 * time.Second

// CreateGrpcClient 创建GRPC客户端
func CreateGrpcClient(ctx context.Context, r registry.Discovery, serviceName string, cfg *conf.Bootstrap, mds ...middleware.Middleware) grpc.ClientConnInterface {

	var options []kratosGrpc.ClientOption

	options = append(options, kratosGrpc.WithDiscovery(r))

	var endpoint string
	if strings.HasPrefix(serviceName, "discovery:///") {
		endpoint = serviceName
	} else {
		endpoint = "discovery:///" + serviceName
	}
	options = append(options, kratosGrpc.WithEndpoint(endpoint))

	options = append(options, initGrpcClientConfig(cfg, mds...)...)

	conn, err := kratosGrpc.DialInsecure(ctx, options...)
	if err != nil {
		log.Fatalf("dial grpc client [%s] failed: %s", serviceName, err.Error())
	}

	return conn
}

func initGrpcClientConfig(cfg *conf.Bootstrap, mds ...middleware.Middleware) []kratosGrpc.ClientOption {
	if cfg.Client == nil || cfg.Client.Grpc == nil {
		return nil
	}

	var options []kratosGrpc.ClientOption

	timeout := defaultTimeout
	if cfg.Client.Grpc.Timeout != nil {
		timeout = cfg.Client.Grpc.Timeout.AsDuration()
	}
	options = append(options, kratosGrpc.WithTimeout(timeout))

	var ms []middleware.Middleware
	ms = append(ms, mds...)
	if cfg.Client.Grpc.Middleware != nil {
		if cfg.Client.Grpc.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Client.Grpc.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Client())
		}
		if cfg.Client.Grpc.Middleware.GetEnableValidate() {
			ms = append(ms, validate.Validator())
		}
		if cfg.Client.Grpc.Middleware.GetEnableMetadata() {
			ms = append(ms, metadata.Client())
		}
	}
	options = append(options, kratosGrpc.WithMiddleware(ms...))

	if cfg.Client.Grpc.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = utils.LoadClientTlsConfig(cfg.Client.Grpc.Tls); err != nil {
			panic(err)
		}

		if tlsCfg != nil {
			options = append(options, kratosGrpc.WithTLSConfig(tlsCfg))
		}
	}

	return options
}

// CreateGrpcServer 创建GRPC服务端
func CreateGrpcServer(cfg *conf.Bootstrap, mds ...middleware.Middleware) *kratosGrpc.Server {
	var options []kratosGrpc.ServerOption

	options = append(options, initGrpcServerConfig(cfg, mds...)...)

	srv := kratosGrpc.NewServer(options...)

	return srv
}

func initGrpcServerConfig(cfg *conf.Bootstrap, mds ...middleware.Middleware) []kratosGrpc.ServerOption {
	if cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil
	}

	var options []kratosGrpc.ServerOption

	var ms []middleware.Middleware
	ms = append(ms, mds...)
	if cfg.Server.Grpc.Middleware != nil {
		if cfg.Server.Grpc.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Grpc.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Grpc.Middleware.GetEnableValidate() {
			ms = append(ms, validate.Validator())
		}
		if cfg.Server.Grpc.Middleware.GetEnableCircuitBreaker() {
		}
		if cfg.Server.Grpc.Middleware.Limiter != nil {
			var limiter ratelimit.Limiter
			switch cfg.Server.Grpc.Middleware.Limiter.GetName() {
			case "bbr":
				limiter = bbr.NewLimiter()
			}
			ms = append(ms, midRateLimit.Server(midRateLimit.WithLimiter(limiter)))
		}
		if cfg.Server.Grpc.Middleware.GetEnableMetadata() {
			ms = append(ms, metadata.Server())
		}
	}
	options = append(options, kratosGrpc.Middleware(ms...))

	if cfg.Server.Grpc.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = utils.LoadServerTlsConfig(cfg.Server.Grpc.Tls); err != nil {
			panic(err)
		}

		if tlsCfg != nil {
			options = append(options, kratosGrpc.TLSConfig(tlsCfg))
		}
	}

	if cfg.Server.Grpc.Network != "" {
		options = append(options, kratosGrpc.Network(cfg.Server.Grpc.Network))
	}
	if cfg.Server.Grpc.Addr != "" {
		options = append(options, kratosGrpc.Address(cfg.Server.Grpc.Addr))
	}
	if cfg.Server.Grpc.Timeout != nil {
		options = append(options, kratosGrpc.Timeout(cfg.Server.Grpc.Timeout.AsDuration()))
	}

	return options
}
