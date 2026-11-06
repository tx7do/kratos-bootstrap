package rpc

import (
	"crypto/tls"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/gorilla/handlers"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"

	kratosRest "github.com/go-kratos/kratos/v2/transport/http"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	midRateLimit "github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/rpc/middleware/validate"
)

// CreateRestServer 创建REST服务端
func CreateRestServer(cfg *conf.Bootstrap, mds ...middleware.Middleware) (*kratosRest.Server, error) {
	options, err := initRestConfig(cfg, mds...)
	if err != nil {
		return nil, err
	}

	ppf := cfg != nil && cfg.Server != nil && cfg.Server.Rest != nil && cfg.Server.Rest.GetEnablePprof()
	if !ppf {
		// pprof 开关关闭时必须前置拦截整个 /debug/ 面：
		// 本包 import 了 net/http/pprof，其 init() 无条件向 http.DefaultServeMux
		// 注册全部 pprof handler；而 kratos mux 未命中业务路由时把请求交给
		// DefaultServeMux（kratos transport/http server.go 的 NotFoundHandler），
		// 若不拦截，无论开关真假 pprof 端点恒可访问。
		// 用 Filter（mux 匹配外层）而非注册 deny 路由，结构上保证不可绕过；
		// 业务 API 无 /debug/ 前缀，整段封禁不影响业务。
		options = append(options, kratosRest.Filter(blockDebugPprof()))
	}

	srv := kratosRest.NewServer(options...)

	if ppf {
		registerHttpPprof(srv)
	}

	return srv, nil
}

// initRestConfig 初始化REST服务配置
func initRestConfig(cfg *conf.Bootstrap, mds ...middleware.Middleware) ([]kratosRest.ServerOption, error) {
	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil, nil
	}

	var options []kratosRest.ServerOption

	if cfg.Server.Rest.Cors != nil {
		corsOptions := []handlers.CORSOption{
			handlers.AllowedHeaders(cfg.Server.Rest.Cors.Headers),
			handlers.AllowedMethods(cfg.Server.Rest.Cors.Methods),
			handlers.AllowedOrigins(cfg.Server.Rest.Cors.Origins),
		}
		// 前端以 HttpOnly Cookie 传 refresh token 时（withCredentials=true），
		// 浏览器要求预检响应携带 Access-Control-Allow-Credentials: true
		if cfg.Server.Rest.Cors.GetAllowCredentials() {
			corsOptions = append(corsOptions, handlers.AllowCredentials())
		}
		options = append(options, kratosRest.Filter(handlers.CORS(corsOptions...)))
	}

	var ms []middleware.Middleware
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
		if cfg.Server.Rest.Middleware.GetEnableMetadata() {
			ms = append(ms, metadata.Server())
		}
	}
	ms = append(ms, mds...)

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

		if tlsCfg, err = loadServerTlsConfig(cfg.Server.Rest.Tls); err != nil {
			return nil, err
		}

		if tlsCfg != nil {
			options = append(options, kratosRest.TLSConfig(tlsCfg))
		}
	}

	return options, nil
}

// registerHttpPprof 注册pprof路由
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

// blockDebugPprof 前置拦截整个 /debug/ 面（pprof 关闭态）。
// kratos Filter 位于 mux 匹配外层，优先于 NotFoundHandler 的 DefaultServeMux
// fallback；net/http/pprof init() 挂在 DefaultServeMux 上的 handler 因此不可达。
func blockDebugPprof() kratosRest.FilterFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/debug" || strings.HasPrefix(r.URL.Path, "/debug/") {
				http.Error(w, "debug endpoints are disabled", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// NewRestWhiteListMatcher 创建REST白名单匹配器
func NewRestWhiteListMatcher() selector.MatchFunc {
	// reuse package-level DefaultWhiteList matcher for REST
	return NewWhiteListMatcher()
}
