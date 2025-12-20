package etcd

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/config"
)

type source struct {
	client  *clientv3.Client
	options *options
}

func New(client *clientv3.Client, opts ...Option) (config.Source, error) {
	o := &options{
		ctx:    context.Background(),
		path:   "",
		prefix: false,
	}

	for _, opt := range opts {
		opt(o)
	}

	if o.path == "" {
		return nil, errors.New("path invalid")
	}

	return &source{
		client:  client,
		options: o,
	}, nil
}

// Load return the config values
func (s *source) Load() ([]*config.KeyValue, error) {
	var opts []clientv3.OpOption
	if s.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}

	rsp, err := s.client.Get(s.options.ctx, s.options.path, opts...)
	if err != nil {
		return nil, wrapConnError("get key", s.options.path, err)
	}

	kvs := make([]*config.KeyValue, 0, len(rsp.Kvs))
	for _, item := range rsp.Kvs {
		k := string(item.Key)
		kvs = append(kvs, &config.KeyValue{
			Key:    k,
			Value:  item.Value,
			Format: strings.TrimPrefix(filepath.Ext(k), "."),
		})
	}
	return kvs, nil
}

// Watch return the watcher
func (s *source) Watch() (config.Watcher, error) {
	w, err := newWatcher(s)
	if err != nil {
		return nil, wrapConnError("create watcher", s.options.path, err)
	}
	return w, nil
}

// isConnError 判断是否为连接/网络相关错误
func isConnError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}
	e := strings.ToLower(err.Error())
	indicators := []string{
		"connection refused", "connection reset", "no available endpoints",
		"transport is closing", "i/o timeout", "timeout", "connection timed out",
		"tls:", "connection reset by peer", "eof",
	}
	for _, sub := range indicators {
		if strings.Contains(e, sub) {
			return true
		}
	}
	return false
}

// wrapConnError 为连接相关错误提供更明确的提示
func wrapConnError(op string, path string, err error) error {
	if err == nil {
		return nil
	}
	if isConnError(err) {
		if path == "" {
			return fmt.Errorf("etcd: %s failed (cannot reach etcd server): %w", op, err)
		}
		return fmt.Errorf("etcd: %s failed for path %s (cannot reach etcd server): %w", op, path, err)
	}
	return err
}
