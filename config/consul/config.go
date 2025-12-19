package consul

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"

	"github.com/hashicorp/consul/api"
)

type source struct {
	client  *api.Client
	options *options
}

func New(client *api.Client, opts ...Option) (config.Source, error) {
	o := &options{
		ctx:  context.Background(),
		path: "",
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
	kv, _, err := s.client.KV().List(s.options.path, nil)
	if err != nil {
		return nil, err
	}

	pathPrefix := s.options.path
	if !strings.HasSuffix(s.options.path, "/") {
		pathPrefix = pathPrefix + "/"
	}
	kvs := make([]*config.KeyValue, 0)
	for _, item := range kv {
		k := strings.TrimPrefix(item.Key, pathPrefix)
		if k == "" {
			continue
		}
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
	return newWatcher(s)
}
