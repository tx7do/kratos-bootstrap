package rpc

import (
	"crypto/tls"
	"net/http/pprof"

	"github.com/gorilla/handlers"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"

	"github.com/go-kratos/kratos/v2/middleware"
	midRateLimit "github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"

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

		if cfg.Server.Rest.Tls.File != nil {
			if tlsCfg, err = utils.LoadServerTlsConfigFile(
				cfg.Server.Rest.Tls.File.GetKeyPath(),
				cfg.Server.Rest.Tls.File.GetCertPath(),
				cfg.Server.Rest.Tls.File.GetCaPath(),
				cfg.Server.Rest.Tls.InsecureSkipVerify,
			); err != nil {
				panic(err)
			}
		}
		if tlsCfg == nil && cfg.Server.Rest.Tls.Config != nil {
			if tlsCfg, err = utils.LoadServerTlsConfig(
				cfg.Server.Rest.Tls.Config.GetKeyPem(),
				cfg.Server.Rest.Tls.Config.GetCertPem(),
				cfg.Server.Rest.Tls.Config.GetCaPem(),
				cfg.Server.Rest.Tls.InsecureSkipVerify,
			); err != nil {
				panic(err)
			}
		}

		if tlsCfg != nil {
			options = append(options, kratosRest.TLSConfig(tlsCfg))
		}
	}

	return options
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
