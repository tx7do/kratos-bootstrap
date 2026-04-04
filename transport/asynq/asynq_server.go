package asynq

import (
	"crypto/tls"

	"github.com/hibiken/asynq"

	tlsUtils "github.com/tx7do/go-utils/tls"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	asynqTransport "github.com/tx7do/kratos-transport/transport/asynq"
)

// NewAsynqServer creates a new Asynq server.
func NewAsynqServer(cfg *conf.Server_Asynq, opts ...asynqTransport.ServerOption) *asynqTransport.Server {
	if cfg == nil {
		return nil
	}

	var o []asynqTransport.ServerOption

	if cfg.GetCodec() != "" {
		o = append(o, asynqTransport.WithCodec(cfg.GetCodec()))
	}

	switch cfg.RedisClientOpts.(type) {
	case *conf.Server_Asynq_RedisOpt:
		redisOpt := cfg.GetRedisOpt()

		opt := &asynq.RedisClientOpt{}

		if redisOpt.GetNetwork() != "" {
			opt.Network = redisOpt.GetNetwork()
		}
		if redisOpt.GetAddr() != "" {
			opt.Addr = redisOpt.GetAddr()
		}
		if redisOpt.GetPassword() != "" {
			opt.Password = redisOpt.GetPassword()
		}
		opt.DB = int(redisOpt.GetDb())
		opt.PoolSize = int(redisOpt.GetPoolSize())

		if cfg.DialTimeout != nil {
			opt.DialTimeout = cfg.GetDialTimeout().AsDuration()
		}
		if cfg.ReadTimeout != nil {
			opt.ReadTimeout = cfg.GetReadTimeout().AsDuration()
		}
		if cfg.WriteTimeout != nil {
			opt.WriteTimeout = cfg.GetWriteTimeout().AsDuration()
		}

		if cfg.Tls != nil {
			if tlsCfg, err := loadClientTlsConfig(cfg.Tls); err != nil {
			} else {
				opt.TLSConfig = tlsCfg
			}
		}

		o = append(o, asynqTransport.WithRedisConnOpt(opt))

	case *conf.Server_Asynq_RedisClusterOpt:
		redisClusterOpt := cfg.GetRedisClusterOpt()

		opt := &asynq.RedisClusterClientOpt{}

		opt.Addrs = redisClusterOpt.GetAddrs()
		opt.MaxRedirects = int(redisClusterOpt.GetMaxRedirects())
		opt.Username = redisClusterOpt.GetUsername()
		opt.Password = redisClusterOpt.GetPassword()

		if cfg.DialTimeout != nil {
			opt.DialTimeout = cfg.GetDialTimeout().AsDuration()
		}
		if cfg.ReadTimeout != nil {
			opt.ReadTimeout = cfg.GetReadTimeout().AsDuration()
		}
		if cfg.WriteTimeout != nil {
			opt.WriteTimeout = cfg.GetWriteTimeout().AsDuration()
		}

		if cfg.Tls != nil {
			if tlsCfg, err := loadClientTlsConfig(cfg.Tls); err != nil {
			} else {
				opt.TLSConfig = tlsCfg
			}
		}

		o = append(o, asynqTransport.WithRedisConnOpt(opt))

	case *conf.Server_Asynq_RedisFailoverOpt:
		redisFailoverOpt := cfg.GetRedisFailoverOpt()

		opt := &asynq.RedisFailoverClientOpt{}

		opt.MasterName = redisFailoverOpt.GetMasterName()
		opt.SentinelAddrs = redisFailoverOpt.GetSentinelAddrs()
		opt.SentinelUsername = redisFailoverOpt.GetSentinelUsername()
		opt.SentinelPassword = redisFailoverOpt.GetSentinelPassword()

		opt.Username = redisFailoverOpt.GetUsername()
		opt.Password = redisFailoverOpt.GetPassword()

		opt.PoolSize = int(redisFailoverOpt.GetPoolSize())

		if cfg.DialTimeout != nil {
			opt.DialTimeout = cfg.GetDialTimeout().AsDuration()
		}
		if cfg.ReadTimeout != nil {
			opt.ReadTimeout = cfg.GetReadTimeout().AsDuration()
		}
		if cfg.WriteTimeout != nil {
			opt.WriteTimeout = cfg.GetWriteTimeout().AsDuration()
		}

		if cfg.Tls != nil {
			if tlsCfg, err := loadClientTlsConfig(cfg.Tls); err != nil {
			} else {
				opt.TLSConfig = tlsCfg
			}
		}

		o = append(o, asynqTransport.WithRedisConnOpt(opt))

	case *conf.Server_Asynq_Uri:
		if cfg.GetUri() != "" {
			o = append(o, asynqTransport.WithRedisURI(cfg.GetUri()))
		}

	default:

	}

	if cfg.GetLocation() != "" {
		o = append(o, asynqTransport.WithLocation(cfg.GetLocation()))
	}
	if cfg.Concurrency != nil {
		o = append(o, asynqTransport.WithConcurrency(cfg.GetConcurrency()))
	}
	if cfg.GroupMaxSize != nil {
		o = append(o, asynqTransport.WithGroupMaxSize(cfg.GetGroupMaxSize()))
	}
	if len(cfg.Queues) != 0 {
		o = append(o, asynqTransport.WithQueues(cfg.GetQueues()))
	}

	if cfg.EnableGracefullyShutdown != nil {
		o = append(o, asynqTransport.WithGracefullyShutdown(cfg.GetEnableGracefullyShutdown()))
	}
	if cfg.EnableStrictPriority != nil {
		o = append(o, asynqTransport.WithStrictPriority(cfg.GetEnableStrictPriority()))
	}

	if cfg.ShutdownTimeout != nil {
		o = append(o, asynqTransport.WithShutdownTimeout(cfg.GetShutdownTimeout().AsDuration()))
	}

	if cfg.TaskCheckInterval != nil {
		o = append(o, asynqTransport.WithTaskCheckInterval(cfg.GetTaskCheckInterval().AsDuration()))
	}
	if cfg.HealthCheckInterval != nil {
		o = append(o, asynqTransport.WithHealthCheckInterval(cfg.GetHealthCheckInterval().AsDuration()))
	}
	if cfg.DelayedTaskCheckInterval != nil {
		o = append(o, asynqTransport.WithDelayedTaskCheckInterval(cfg.GetDelayedTaskCheckInterval().AsDuration()))
	}
	if cfg.GroupGracePeriod != nil {
		o = append(o, asynqTransport.WithGroupGracePeriod(cfg.GetGroupGracePeriod().AsDuration()))
	}
	if cfg.GroupMaxDelay != nil {
		o = append(o, asynqTransport.WithGroupMaxDelay(cfg.GetGroupMaxDelay().AsDuration()))
	}
	if cfg.JanitorInterval != nil {
		o = append(o, asynqTransport.WithJanitorInterval(cfg.GetJanitorInterval().AsDuration()))
	}
	if cfg.JanitorBatchSize != nil {
		o = append(o, asynqTransport.WithJanitorBatchSize(cfg.GetJanitorBatchSize()))
	}

	if opts != nil {
		o = append(o, opts...)
	}

	srv := asynqTransport.NewServer(o...)

	return srv
}

func loadClientTlsConfig(cfg *conf.TLS) (*tls.Config, error) {
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

	return tlsCfg, err
}
