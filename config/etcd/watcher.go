package etcd

import (
	"context"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/config"
)

type watcher struct {
	source *source
	ch     clientv3.WatchChan

	ctx    context.Context
	cancel context.CancelFunc
}

func newWatcher(s *source) (*watcher, error) {
	// 短超时探测 etcd 可达性
	probeCtx, probeCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer probeCancel()
	if _, err := s.client.Get(probeCtx, s.options.path, clientv3.WithLimit(1)); err != nil {
		return nil, wrapConnError("create watcher", s.options.path, err)
	}

	// 创建 watcher
	ctx, cancel := context.WithCancel(context.Background())
	w := &watcher{
		source: s,
		ctx:    ctx,
		cancel: cancel,
	}

	var opts []clientv3.OpOption
	if s.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}
	w.ch = s.client.Watch(s.options.ctx, s.options.path, opts...)

	return w, nil
}

func (w *watcher) Next() ([]*config.KeyValue, error) {
	select {
	case resp := <-w.ch:
		if err := resp.Err(); err != nil {
			return nil, err
		}
		return w.source.Load()
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	}
}

func (w *watcher) Stop() error {
	w.cancel()
	return nil
}
