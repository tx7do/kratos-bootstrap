package asynq

import (
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-transport/transport/asynq"
)

// NewAsynqServer creates a new asynq server.
func NewAsynqServer(cfg *conf.Server, opts ...asynq.ServerOption) *asynq.Server {
	if cfg == nil || cfg.Asynq == nil {
		return nil
	}

	var o []asynq.ServerOption

	if cfg.Asynq.GetCodec() != "" {
		o = append(o, asynq.WithCodec(cfg.Asynq.GetCodec()))
	}

	if cfg.Asynq.GetUri() != "" {
		o = append(o, asynq.WithRedisURI(cfg.Asynq.GetUri()))
	}
	if cfg.Asynq.Db != 0 {
		o = append(o, asynq.WithRedisDB(int(cfg.Asynq.GetDb())))
	}
	if cfg.Asynq.PoolSize != 0 {
		o = append(o, asynq.WithRedisPoolSize(int(cfg.Asynq.GetPoolSize())))
	}

	if cfg.Asynq.GetLocation() != "" {
		o = append(o, asynq.WithLocation(cfg.Asynq.GetLocation()))
	}
	if cfg.Asynq.Concurrency != 0 {
		o = append(o, asynq.WithConcurrency(int(cfg.Asynq.GetConcurrency())))
	}
	if cfg.Asynq.GroupMaxSize != 0 {
		o = append(o, asynq.WithGroupMaxSize(int(cfg.Asynq.GetGroupMaxSize())))
	}
	if len(cfg.Asynq.Queues) != 0 {
		var queues = make(map[string]int)
		for k, v := range cfg.Asynq.GetQueues() {
			queues[k] = int(v)
		}
		o = append(o, asynq.WithQueues(queues))
	}

	if cfg.Asynq.ShutdownTimeout != nil {
		o = append(o, asynq.WithShutdownTimeout(cfg.Asynq.GetShutdownTimeout().AsDuration()))
	}
	if cfg.Asynq.DialTimeout != nil {
		o = append(o, asynq.WithDialTimeout(cfg.Asynq.GetDialTimeout().AsDuration()))
	}
	if cfg.Asynq.ReadTimeout != nil {
		o = append(o, asynq.WithReadTimeout(cfg.Asynq.GetReadTimeout().AsDuration()))
	}
	if cfg.Asynq.WriteTimeout != nil {
		o = append(o, asynq.WithWriteTimeout(cfg.Asynq.GetWriteTimeout().AsDuration()))
	}

	if cfg.Asynq.HealthCheckInterval != nil {
		o = append(o, asynq.WithHealthCheckInterval(cfg.Asynq.GetHealthCheckInterval().AsDuration()))
	}
	if cfg.Asynq.DelayedTaskCheckInterval != nil {
		o = append(o, asynq.WithDelayedTaskCheckInterval(cfg.Asynq.GetDelayedTaskCheckInterval().AsDuration()))
	}
	if cfg.Asynq.GroupMaxDelay != nil {
		o = append(o, asynq.WithGroupMaxDelay(cfg.Asynq.GetGroupMaxDelay().AsDuration()))
	}
	if cfg.Asynq.GroupGracePeriod != nil {
		o = append(o, asynq.WithGroupGracePeriod(cfg.Asynq.GetGroupGracePeriod().AsDuration()))
	}

	o = append(o,
		asynq.WithGracefullyShutdown(cfg.Asynq.GetEnableGracefullyShutdown()),
		asynq.WithShutdownTimeout(cfg.Asynq.GetShutdownTimeout().AsDuration()),
	)

	if opts != nil {
		o = append(o, opts...)
	}

	srv := asynq.NewServer(o...)

	return srv
}
