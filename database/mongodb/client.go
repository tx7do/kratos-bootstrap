package mongodb

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// NewMongoClient 创建MongoDB客户端
func NewMongoClient(ctx context.Context, cfg *conf.Bootstrap, l *log.Helper) *mongo.Client {
	if cfg.Data == nil || cfg.Data.Mongodb == nil {
		l.Warn("Mongodb config is nil")
		return nil
	}

	var opts []*options.ClientOptions

	uri := fmt.Sprintf("mongodb://%s:%s@%s",
		cfg.Data.Mongodb.Username, cfg.Data.Mongodb.Password, cfg.Data.Mongodb.Address,
	)
	opts = append(opts, options.Client().ApplyURI(uri))

	cli, err := mongo.Connect(ctx, opts...)
	if err != nil {
		l.Fatalf("failed opening connection to mongodb: %v", err)
		return nil
	}

	return cli
}
