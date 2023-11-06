package bootstrap

import (
	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"

	conf "github.com/tx7do/kratos-bootstrap/gen/api/go/conf/v1"
)

// NewRedisClient 创建Redis客户端
func NewRedisClient(conf *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.GetRedis().GetAddr(),
		Password:     conf.GetRedis().GetPassword(),
		DB:           int(conf.GetRedis().GetDb()),
		DialTimeout:  conf.GetRedis().GetDialTimeout().AsDuration(),
		WriteTimeout: conf.GetRedis().GetWriteTimeout().AsDuration(),
		ReadTimeout:  conf.GetRedis().GetReadTimeout().AsDuration(),
	})
	if rdb == nil {
		log.Fatalf("failed opening connection to redis")
		return nil
	}
	rdb.AddHook(redisotel.NewTracingHook())

	return rdb
}
