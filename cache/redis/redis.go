package redis

import (
	"crypto/tls"
	"errors"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	tlsUtils "github.com/tx7do/go-utils/tls"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewClient create go-redis client (standalone mode only).
func NewClient(conf *conf.Data, logger *log.Helper) (rdb *redis.Client) {
	in := conf.GetRedis()

	opts := &redis.Options{
		Addr:         in.GetAddr(),
		Username:     in.GetUsername(),
		Password:     in.GetPassword(),
		DB:           int(in.GetDb()),
		Network:      in.GetNetwork(),
		DialTimeout:  in.GetDialTimeout().AsDuration(),
		WriteTimeout: in.GetWriteTimeout().AsDuration(),
		ReadTimeout:  in.GetReadTimeout().AsDuration(),
		PoolTimeout:  in.GetPoolTimeout().AsDuration(),
		ConnMaxIdleTime: in.GetConnMaxIdleTime().AsDuration(),
		ConnMaxLifetime: in.GetConnMaxLifetime().AsDuration(),
		MinRetryBackoff: in.GetMinRetryBackoff().AsDuration(),
		MaxRetryBackoff: in.GetMaxRetryBackoff().AsDuration(),
	}

	if in.PoolSize != nil {
		opts.PoolSize = int(in.GetPoolSize())
	}
	if in.MinIdleConns != nil {
		opts.MinIdleConns = int(in.GetMinIdleConns())
	}
	if in.MaxRetries != nil {
		opts.MaxRetries = int(in.GetMaxRetries())
	}

	if tlsCfg, err := loadTlsConfig(in.GetTls()); err != nil {
		logger.Errorf("failed load tls config: %s", err.Error())
		return nil
	} else if tlsCfg != nil {
		opts.TLSConfig = tlsCfg
	}

	if rdb = redis.NewClient(opts); rdb == nil {
		logger.Errorf("failed opening connection to redis")
		return nil
	}

	instrument(in, rdb, logger)

	return rdb
}

// NewUniversalClient create go-redis client by mode: standalone、cluster、sentinel.
func NewUniversalClient(conf *conf.Data, logger *log.Helper) (redis.UniversalClient, error) {
	in := conf.GetRedis()

	switch strings.ToLower(in.GetMode()) {
	case "", "standalone":
		return NewClient(conf, logger), nil

	case "cluster":
		opts := &redis.ClusterOptions{
			Addrs:          in.GetAddrs(),
			Username:       in.GetUsername(),
			Password:       in.GetPassword(),
			DialTimeout:    in.GetDialTimeout().AsDuration(),
			ReadTimeout:    in.GetReadTimeout().AsDuration(),
			WriteTimeout:   in.GetWriteTimeout().AsDuration(),
			PoolTimeout:    in.GetPoolTimeout().AsDuration(),
			ConnMaxIdleTime: in.GetConnMaxIdleTime().AsDuration(),
			ConnMaxLifetime: in.GetConnMaxLifetime().AsDuration(),
			MinRetryBackoff: in.GetMinRetryBackoff().AsDuration(),
			MaxRetryBackoff: in.GetMaxRetryBackoff().AsDuration(),
			RouteByLatency: in.GetRouteByLatency(),
			RouteRandomly:  in.GetRouteRandomly(),
		}
		if in.PoolSize != nil {
			opts.PoolSize = int(in.GetPoolSize())
		}
		if in.MinIdleConns != nil {
			opts.MinIdleConns = int(in.GetMinIdleConns())
		}
		if in.MaxRetries != nil {
			opts.MaxRetries = int(in.GetMaxRetries())
		}
		if len(opts.Addrs) == 0 && in.GetAddr() != "" {
			opts.Addrs = strings.Split(in.GetAddr(), ",")
		}
		tlsCfg, err := loadTlsConfig(in.GetTls())
		if err != nil {
			return nil, err
		}
		opts.TLSConfig = tlsCfg

		cli := redis.NewClusterClient(opts)
		instrument(in, cli, logger)
		return cli, nil

	case "sentinel":
		opts := &redis.FailoverOptions{
			MasterName:       in.GetMasterName(),
			SentinelAddrs:    in.GetAddrs(),
			SentinelUsername: in.GetSentinelUsername(),
			SentinelPassword: in.GetSentinelPassword(),
			Username:         in.GetUsername(),
			Password:         in.GetPassword(),
			DB:               int(in.GetDb()),
			DialTimeout:      in.GetDialTimeout().AsDuration(),
			ReadTimeout:      in.GetReadTimeout().AsDuration(),
			WriteTimeout:     in.GetWriteTimeout().AsDuration(),
			PoolTimeout:      in.GetPoolTimeout().AsDuration(),
			ConnMaxIdleTime:  in.GetConnMaxIdleTime().AsDuration(),
			ConnMaxLifetime:  in.GetConnMaxLifetime().AsDuration(),
			MinRetryBackoff:  in.GetMinRetryBackoff().AsDuration(),
			MaxRetryBackoff:  in.GetMaxRetryBackoff().AsDuration(),
			RouteByLatency:   in.GetRouteByLatency(),
			RouteRandomly:    in.GetRouteRandomly(),
			ReplicaOnly:      in.GetReplicaOnly(),
		}
		if in.PoolSize != nil {
			opts.PoolSize = int(in.GetPoolSize())
		}
		if in.MinIdleConns != nil {
			opts.MinIdleConns = int(in.GetMinIdleConns())
		}
		if in.MaxRetries != nil {
			opts.MaxRetries = int(in.GetMaxRetries())
		}
		if len(opts.SentinelAddrs) == 0 && in.GetAddr() != "" {
			opts.SentinelAddrs = strings.Split(in.GetAddr(), ",")
		}
		tlsCfg, err := loadTlsConfig(in.GetTls())
		if err != nil {
			return nil, err
		}
		opts.TLSConfig = tlsCfg

		cli := redis.NewFailoverClient(opts)
		instrument(in, cli, logger)
		return cli, nil
	}

	return nil, errors.New("redis: unknown mode: " + in.GetMode())
}

func instrument(in *conf.Data_Redis, cli redis.UniversalClient, logger *log.Helper) {
	// open tracing instrumentation.
	if in.GetEnableTracing() {
		if err := redisotel.InstrumentTracing(cli); err != nil {
			logger.Errorf("failed open tracing: %s", err.Error())
		}
	}

	// open metrics instrumentation.
	if in.GetEnableMetrics() {
		if err := redisotel.InstrumentMetrics(cli); err != nil {
			logger.Errorf("failed open metrics: %s", err.Error())
		}
	}
}

func loadTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	var tlsCfg *tls.Config
	var err error

	if cfg.File != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigFile(
			cfg.File.GetKeyPath(),
			cfg.File.GetCertPath(),
			cfg.File.GetCaPath(),
		); err != nil {
			return nil, err
		}
	} else if cfg.Config != nil {
		if tlsCfg, err = tlsUtils.LoadClientTlsConfigString(
			cfg.Config.GetKeyPem(),
			cfg.Config.GetCertPem(),
			cfg.Config.GetCaPem(),
		); err != nil {
			return nil, err
		}
	}

	if tlsCfg != nil && cfg.GetInsecureSkipVerify() {
		tlsCfg.InsecureSkipVerify = true
	}

	return tlsCfg, nil
}
