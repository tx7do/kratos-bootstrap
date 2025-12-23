package asynq

import (
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-transport/transport/asynq"
)

// NewAsynqServer creates a new Asynq server.
func NewAsynqServer(cfg *conf.Server_Asynq, opts ...asynq.ServerOption) *asynq.Server {
	if cfg == nil {
		return nil
	}

	var o []asynq.ServerOption

	if cfg.GetCodec() != "" {
		o = append(o, asynq.WithCodec(cfg.GetCodec()))
	}

	if cfg.GetUri() != "" {
		o = append(o, asynq.WithRedisURI(cfg.GetUri()))
	}
	if cfg.Db != 0 {
		o = append(o, asynq.WithRedisDB(int(cfg.GetDb())))
	}
	if cfg.PoolSize != 0 {
		o = append(o, asynq.WithRedisPoolSize(int(cfg.GetPoolSize())))
	}

	if cfg.GetLocation() != "" {
		o = append(o, asynq.WithLocation(cfg.GetLocation()))
	}
	if cfg.Concurrency != 0 {
		o = append(o, asynq.WithConcurrency(int(cfg.GetConcurrency())))
	}
	if cfg.GroupMaxSize != 0 {
		o = append(o, asynq.WithGroupMaxSize(int(cfg.GetGroupMaxSize())))
	}
	if len(cfg.Queues) != 0 {
		var queues = make(map[string]int)
		for k, v := range cfg.GetQueues() {
			queues[k] = int(v)
		}
		o = append(o, asynq.WithQueues(queues))
	}

	if cfg.ShutdownTimeout != nil {
		o = append(o, asynq.WithShutdownTimeout(cfg.GetShutdownTimeout().AsDuration()))
	}
	if cfg.DialTimeout != nil {
		o = append(o, asynq.WithDialTimeout(cfg.GetDialTimeout().AsDuration()))
	}
	if cfg.ReadTimeout != nil {
		o = append(o, asynq.WithReadTimeout(cfg.GetReadTimeout().AsDuration()))
	}
	if cfg.WriteTimeout != nil {
		o = append(o, asynq.WithWriteTimeout(cfg.GetWriteTimeout().AsDuration()))
	}

	if cfg.HealthCheckInterval != nil {
		o = append(o, asynq.WithHealthCheckInterval(cfg.GetHealthCheckInterval().AsDuration()))
	}
	if cfg.DelayedTaskCheckInterval != nil {
		o = append(o, asynq.WithDelayedTaskCheckInterval(cfg.GetDelayedTaskCheckInterval().AsDuration()))
	}
	if cfg.GroupMaxDelay != nil {
		o = append(o, asynq.WithGroupMaxDelay(cfg.GetGroupMaxDelay().AsDuration()))
	}
	if cfg.GroupGracePeriod != nil {
		o = append(o, asynq.WithGroupGracePeriod(cfg.GetGroupGracePeriod().AsDuration()))
	}

	o = append(o,
		asynq.WithGracefullyShutdown(cfg.GetEnableGracefullyShutdown()),
		asynq.WithShutdownTimeout(cfg.GetShutdownTimeout().AsDuration()),
	)

	if opts != nil {
		o = append(o, opts...)
	}

	srv := asynq.NewServer(o...)

	return srv
}
