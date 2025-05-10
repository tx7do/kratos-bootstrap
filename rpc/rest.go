package rpc

import (
	"crypto/tls"
	"net/http/pprof"

	"github.com/gorilla/handlers"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	midRateLimit "github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	kratosRest "github.com/go-kratos/kratos/v2/transport/http"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/utils"
)

// CreateRestServer 创建REST服务端
func CreateRestServer(cfg *conf.Bootstrap, mds ...middleware.Middleware) *kratosRest.Server {
	var options []kratosRest.ServerOption

	options = append(options, initRestConfig(cfg, mds...)...)

	srv := kratosRest.NewServer(options...)

	if cfg.Server != nil && cfg.Server.Rest != nil && cfg.Server.Rest.GetEnablePprof() {
		registerHttpPprof(srv)
	}

	return srv
}

func initRestConfig(cfg *conf.Bootstrap, mds ...middleware.Middleware) []kratosRest.ServerOption {
	if cfg.Server == nil || cfg.Server.Rest == nil {
		return nil
	}

	var options []kratosRest.ServerOption

	if cfg.Server.Rest.Cors != nil {
		options = append(options, kratosRest.Filter(handlers.CORS(
			handlers.AllowedHeaders(cfg.Server.Rest.Cors.Headers),
			handlers.AllowedMethods(cfg.Server.Rest.Cors.Methods),
			handlers.AllowedOrigins(cfg.Server.Rest.Cors.Origins),
		)))
	}

	var ms []middleware.Middleware
	ms = append(ms, mds...)
	if cfg.Server.Rest.Middleware != nil {
		if cfg.Server.Rest.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Rest.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Rest.Middleware.GetEnableValidate() {
			ms = append(ms, validate.ProtoValidate())
		}
		if cfg.Server.Rest.Middleware.GetEnableCircuitBreaker() {
		}
		if cfg.Server.Rest.Middleware.Limiter != nil {
			var limiter ratelimit.Limiter
			switch cfg.Server.Rest.Middleware.Limiter.GetName() {
			case "bbr":
				limiter = bbr.NewLimiter()
			}
			ms = append(ms, midRateLimit.Server(midRateLimit.WithLimiter(limiter)))
		}
		if cfg.Server.Grpc.Middleware.GetEnableMetadata() {
			ms = append(ms, metadata.Server())
		}
	}
	options = append(options, kratosRest.Middleware(ms...))

	if cfg.Server.Rest.Network != "" {
		options = append(options, kratosRest.Network(cfg.Server.Rest.Network))
	}
	if cfg.Server.Rest.Addr != "" {
		options = append(options, kratosRest.Address(cfg.Server.Rest.Addr))
	}
	if cfg.Server.Rest.Timeout != nil {
		options = append(options, kratosRest.Timeout(cfg.Server.Rest.Timeout.AsDuration()))
	}

	if cfg.Server.Rest.Tls != nil {
		var tlsCfg *tls.Config
		var err error

		if tlsCfg, err = utils.LoadServerTlsConfig(cfg.Server.Rest.Tls); err != nil {
			panic(err)
		}

		if tlsCfg != nil {
			options = append(options, kratosRest.TLSConfig(tlsCfg))
		}
	}

	return options
}

func registerHttpPprof(s *kratosRest.Server) {
	s.HandleFunc("/debug/pprof", pprof.Index)

	s.HandleFunc("/debug/cmdline", pprof.Cmdline)
	s.HandleFunc("/debug/profile", pprof.Profile)
	s.HandleFunc("/debug/symbol", pprof.Symbol)
	s.HandleFunc("/debug/trace", pprof.Trace)

	s.HandleFunc("/debug/allocs", pprof.Handler("allocs").ServeHTTP)
	s.HandleFunc("/debug/block", pprof.Handler("block").ServeHTTP)
	s.HandleFunc("/debug/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.HandleFunc("/debug/heap", pprof.Handler("heap").ServeHTTP)
	s.HandleFunc("/debug/mutex", pprof.Handler("mutex").ServeHTTP)
	s.HandleFunc("/debug/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}
