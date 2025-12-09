package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	mongoV2 "go.mongodb.org/mongo-driver/v2/mongo"
	optionsV2 "go.mongodb.org/mongo-driver/v2/mongo/options"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	mongodbCrud "github.com/tx7do/go-crud/mongodb"
)

type Client struct {
	log *log.Helper

	cli      *mongoV2.Client
	database string
	timeout  time.Duration
}

func NewClient(logger log.Logger, cfg *conf.Bootstrap) (*mongodbCrud.Client, error) {
	if cfg.Data == nil || cfg.Data.Mongodb == nil {
		return nil, errors.New("mongodb config is nil")
	}

	var options []mongodbCrud.Option

	if logger != nil {
		options = append(options, mongodbCrud.WithLogger(logger))
	}

	if cfg.Data.Mongodb.GetUri() != "" {
		options = append(options, mongodbCrud.WithURI(cfg.Data.Mongodb.GetUri()))
	}
	if cfg.Data.Mongodb.GetDatabase() != "" {
		options = append(options, mongodbCrud.WithDatabase(cfg.Data.Mongodb.GetDatabase()))
	}
	if cfg.Data.Mongodb.Timeout != nil {
		options = append(options, mongodbCrud.WithTimeout(cfg.Data.Mongodb.GetTimeout().AsDuration()))
	}
	if cfg.Data.Mongodb.ConnectTimeout != nil {
		options = append(options, mongodbCrud.WithConnectTimeout(cfg.Data.Mongodb.GetConnectTimeout().AsDuration()))
	}
	if cfg.Data.Mongodb.ServerSelectionTimeout != nil {
		options = append(options, mongodbCrud.WithServerSelectionTimeout(cfg.Data.Mongodb.GetServerSelectionTimeout().AsDuration()))
	}
	if cfg.Data.Mongodb.HeartbeatInterval != nil {
		options = append(options, mongodbCrud.WithHeartbeatInterval(cfg.Data.Mongodb.GetHeartbeatInterval().AsDuration()))
	}
	if cfg.Data.Mongodb.LocalThreshold != nil {
		options = append(options, mongodbCrud.WithLocalThreshold(cfg.Data.Mongodb.GetLocalThreshold().AsDuration()))
	}
	if cfg.Data.Mongodb.MaxConnIdleTime != nil {
		options = append(options, mongodbCrud.WithMaxConnIdleTime(cfg.Data.Mongodb.GetMaxConnIdleTime().AsDuration()))
	}
	if cfg.Data.Mongodb.GetUsername() != "" && cfg.Data.Mongodb.GetPassword() != "" {
		options = append(options, mongodbCrud.WithCredentials(cfg.Data.Mongodb.GetUsername(), cfg.Data.Mongodb.GetPassword()))
	}

	options = append(options, mongodbCrud.WithBSONOptions(&optionsV2.BSONOptions{
		UseJSONStructTags: true, // 使用JSON结构标签
	}))

	return mongodbCrud.NewClient(options...)
}

// createMongodbClient 创建MongoDB客户端
func (c *Client) createMongodbClient(cfg *conf.Bootstrap) error {

	var opts []*optionsV2.ClientOptions

	if cfg.Data.Mongodb.GetUri() != "" {
		opts = append(opts, optionsV2.Client().ApplyURI(cfg.Data.Mongodb.GetUri()))
	}
	if cfg.Data.Mongodb.GetUsername() != "" && cfg.Data.Mongodb.GetPassword() != "" {
		credential := optionsV2.Credential{
			Username: cfg.Data.Mongodb.GetUsername(),
			Password: cfg.Data.Mongodb.GetPassword(),
		}

		if cfg.Data.Mongodb.GetPassword() != "" {
			credential.PasswordSet = true
		}

		opts = append(opts, optionsV2.Client().SetAuth(credential))
	}
	if cfg.Data.Mongodb.ConnectTimeout != nil {
		opts = append(opts, optionsV2.Client().SetConnectTimeout(cfg.Data.Mongodb.GetConnectTimeout().AsDuration()))
	}
	if cfg.Data.Mongodb.ServerSelectionTimeout != nil {
		opts = append(opts, optionsV2.Client().SetServerSelectionTimeout(cfg.Data.Mongodb.GetServerSelectionTimeout().AsDuration()))
	}
	if cfg.Data.Mongodb.Timeout != nil {
		opts = append(opts, optionsV2.Client().SetTimeout(cfg.Data.Mongodb.GetTimeout().AsDuration()))
	}

	opts = append(opts, optionsV2.Client().SetBSONOptions(&optionsV2.BSONOptions{
		UseJSONStructTags: true, // 使用JSON结构标签
	}))

	cli, err := mongoV2.Connect(opts...)
	if err != nil {
		c.log.Errorf("failed to create mongodb client: %v", err)
		return err
	}

	c.database = cfg.Data.Mongodb.GetDatabase()
	if cfg.Data.Mongodb.GetTimeout() != nil {
		c.timeout = cfg.Data.Mongodb.GetTimeout().AsDuration()
	} else {
		c.timeout = 10 * time.Second // 默认超时时间
	}

	c.cli = cli

	return nil
}

// Close 关闭MongoDB客户端
func (c *Client) Close() {
	if c.cli == nil {
		c.log.Warn("mongodb client is already closed or not initialized")
		return
	}

	if err := c.cli.Disconnect(context.Background()); err != nil {
		c.log.Errorf("failed to disconnect mongodb client: %v", err)
	} else {
		c.log.Info("mongodb client disconnected successfully")
	}
}

// CheckConnect 检查MongoDB连接状态
func (c *Client) CheckConnect() {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	if err := c.cli.Ping(ctx, nil); err != nil {
		c.log.Errorf("failed to ping mongodb: %v", err)
	} else {
		c.log.Info("mongodb client is connected")
	}
}

// InsertOne 插入单个文档
func (c *Client) InsertOne(ctx context.Context, collection string, document interface{}) (*mongoV2.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.cli.Database(c.database).Collection(collection).InsertOne(ctx, document)
}

// InsertMany 插入多个文档
func (c *Client) InsertMany(ctx context.Context, collection string, documents []interface{}) (*mongoV2.InsertManyResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.cli.Database(c.database).Collection(collection).InsertMany(ctx, documents)
}

// FindOne 查询单个文档
func (c *Client) FindOne(ctx context.Context, collection string, filter interface{}, result interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.cli.Database(c.database).Collection(collection).FindOne(ctx, filter).Decode(result)
}

// Find 查询多个文档
func (c *Client) Find(ctx context.Context, collection string, filter interface{}, results interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cursor, err := c.cli.Database(c.database).Collection(collection).Find(ctx, filter)
	if err != nil {
		c.log.Errorf("failed to find documents in collection %s: %v", collection, err)
		return err
	}
	defer func(cursor *mongoV2.Cursor, ctx context.Context) {
		if err = cursor.Close(ctx); err != nil {
			c.log.Errorf("failed to close cursor: %v", err)
		}
	}(cursor, ctx)

	return cursor.All(ctx, results)
}

// UpdateOne 更新单个文档
func (c *Client) UpdateOne(ctx context.Context, collection string, filter, update interface{}) (*mongoV2.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.cli.Database(c.database).Collection(collection).UpdateOne(ctx, filter, update)
}

// DeleteOne 删除单个文档
func (c *Client) DeleteOne(ctx context.Context, collection string, filter interface{}) (*mongoV2.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.cli.Database(c.database).Collection(collection).DeleteOne(ctx, filter)
}
