package rpc

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"google.golang.org/grpc"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	midRateLimit "github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"github.com/tx7do/kratos-bootstrap/rpc/middleware/validate"

	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

const defaultTimeout = 5 * time.Second

// CreateGrpcClient 创建GRPC客户端
func CreateGrpcClient(ctx context.Context, r registry.Discovery, serviceName string, cfg *conf.Bootstrap, mds ...middleware.Middleware) (grpc.ClientConnInterface, error) {
	var options []kratosGrpc.ClientOption

	options = append(options, kratosGrpc.WithDiscovery(r))

	var endpoint string
	if strings.HasPrefix(serviceName, "discovery:///") {
		endpoint = serviceName
	} else {
		endpoint = "discovery:///" + serviceName
	}
	options = append(options, kratosGrpc.WithEndpoint(endpoint))

	cfgs, err := initGrpcClientConfig(cfg, mds...)
	if err != nil {
		log.Fatalf("init grpc client config failed: %s", err.Error())
		return nil, err
	}

	options = append(options, cfgs...)

	conn, err := kratosGrpc.DialInsecure(ctx, options...)
	if err != nil {
		log.Fatalf("dial grpc client [%s] failed: %s", serviceName, err.Error())
	}

	return conn, nil
}

func initGrpcClientConfig(cfg *conf.Bootstrap, mds ...middleware.Middleware) ([]kratosGrpc.ClientOption, error) {
	if cfg.Client == nil || cfg.Client.Grpc == nil {
		return nil, nil
	}

	var options []kratosGrpc.ClientOption

	timeout := defaultTimeout
	if cfg.Client.Grpc.Timeout != nil {
		timeout = cfg.Client.Grpc.Timeout.AsDuration()
	}
	options = append(options, kratosGrpc.WithTimeout(timeout))

	var ms []middleware.Middleware
	if cfg.Client.Grpc.Middleware != nil {
		if cfg.Client.Grpc.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Client.Grpc.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Client())
		}
		if cfg.Client.Grpc.Middleware.GetEnableValidate() {
			ms = append(ms, validate.ProtoValidate())
		}
		if cfg.Client.Grpc.Middleware.GetEnableMetadata() {
			ms = append(ms, metadata.Client())
		}
	}
	ms = append(ms, mds...)

	options = append(options, kratosGrpc.WithMiddleware(ms...))

	if cfg.Client.Grpc.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = loadClientTlsConfig(cfg.Client.Grpc.Tls); err != nil {
			return nil, err
		}

		if tlsCfg != nil {
			options = append(options, kratosGrpc.WithTLSConfig(tlsCfg))
		}
	}

	return options, nil
}

// CreateGrpcServer 创建GRPC服务端
func CreateGrpcServer(cfg *conf.Bootstrap, mds ...middleware.Middleware) (*kratosGrpc.Server, error) {
	var options []kratosGrpc.ServerOption

	cfgs, err := initGrpcServerConfig(cfg, mds...)
	if err != nil {
		log.Fatalf("init grpc server config failed: %s", err.Error())
		return nil, err
	}

	options = append(options, cfgs...)

	srv := kratosGrpc.NewServer(options...)

	return srv, nil
}

func initGrpcServerConfig(cfg *conf.Bootstrap, mds ...middleware.Middleware) ([]kratosGrpc.ServerOption, error) {
	if cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil, nil
	}

	var options []kratosGrpc.ServerOption

	var ms []middleware.Middleware
	if cfg.Server.Grpc.Middleware != nil {
		if cfg.Server.Grpc.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Grpc.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Grpc.Middleware.GetEnableValidate() {
			ms = append(ms, validate.ProtoValidate())
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
	ms = append(ms, mds...)

	options = append(options, kratosGrpc.Middleware(ms...))

	if cfg.Server.Grpc.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = loadServerTlsConfig(cfg.Server.Grpc.Tls); err != nil {
			return nil, err
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

	return options, nil
}

func NewGrpcWhiteListMatcher(whiteList *WhiteList) selector.MatchFunc {
	return func(ctx context.Context, operation string) bool {
		if operation == "" {
			return true
		}
		op := normalizeOp(operation)
		whiteList.mu.RLock()
		defer whiteList.mu.RUnlock()
		// skip middleware when whitelisted
		return !whiteList.isWhitelistedLocked(op)
	}
}
