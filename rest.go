package bootstrap

import (
	"net/http/pprof"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"

	"github.com/go-kratos/kratos/v2/middleware"
	midRateLimit "github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"

	kratosRest "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/gorilla/handlers"

	conf "github.com/tx7do/kratos-bootstrap/gen/api/go/conf/v1"
)

// CreateRestServer 创建REST服务端
func CreateRestServer(cfg *conf.Bootstrap, m ...middleware.Middleware) *kratosRest.Server {
	var opts = []kratosRest.ServerOption{
		kratosRest.Filter(handlers.CORS(
			handlers.AllowedHeaders(cfg.Server.Rest.Cors.Headers),
			handlers.AllowedMethods(cfg.Server.Rest.Cors.Methods),
			handlers.AllowedOrigins(cfg.Server.Rest.Cors.Origins),
		)),
	}

	var ms []middleware.Middleware
	if cfg.Server != nil && cfg.Server.Rest != nil && cfg.Server.Rest.Middleware != nil {
		if cfg.Server.Rest.Middleware.GetEnableRecovery() {
			ms = append(ms, recovery.Recovery())
		}
		if cfg.Server.Rest.Middleware.GetEnableTracing() {
			ms = append(ms, tracing.Server())
		}
		if cfg.Server.Rest.Middleware.GetEnableValidate() {
			ms = append(ms, validate.Validator())
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
	}
	ms = append(ms, m...)
	opts = append(opts, kratosRest.Middleware(ms...))

	if cfg.Server.Rest.Network != "" {
		opts = append(opts, kratosRest.Network(cfg.Server.Rest.Network))
	}
	if cfg.Server.Rest.Addr != "" {
		opts = append(opts, kratosRest.Address(cfg.Server.Rest.Addr))
	}
	if cfg.Server.Rest.Timeout != nil {
		opts = append(opts, kratosRest.Timeout(cfg.Server.Rest.Timeout.AsDuration()))
	}

	srv := kratosRest.NewServer(opts...)

	if cfg.Server.Rest.GetEnablePprof() {
		registerHttpPprof(srv)
	}

	return srv
}

func registerHttpPprof(s *kratosRest.Server) {
	s.HandleFunc("/debug/pprof", pprof.Index)
	s.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.HandleFunc("/debug/pprof/trace", pprof.Trace)

	s.HandleFunc("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	s.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	s.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	s.HandleFunc("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
	s.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}
