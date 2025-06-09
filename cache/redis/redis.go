package redis

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewClient create go-redis client
func NewClient(conf *conf.Data, logger *log.Helper) (rdb *redis.Client) {
	if rdb = redis.NewClient(&redis.Options{
		Addr:         conf.GetRedis().GetAddr(),
		Password:     conf.GetRedis().GetPassword(),
		DB:           int(conf.GetRedis().GetDb()),
		DialTimeout:  conf.GetRedis().GetDialTimeout().AsDuration(),
		WriteTimeout: conf.GetRedis().GetWriteTimeout().AsDuration(),
		ReadTimeout:  conf.GetRedis().GetReadTimeout().AsDuration(),
	}); rdb == nil {
		logger.Fatalf("failed opening connection to redis")
		return nil
	}

	// open tracing instrumentation.
	if conf.GetRedis().GetEnableTracing() {
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			logger.Fatalf("failed open tracing: %s", err.Error())
			return nil
		}
	}

	// open metrics instrumentation.
	if conf.GetRedis().GetEnableMetrics() {
		if err := redisotel.InstrumentMetrics(rdb); err != nil {
			logger.Fatalf("failed open metrics: %s", err.Error())
			return nil
		}
	}

	return rdb
}
