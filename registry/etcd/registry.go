package etcd

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/registry"
)

var (
	_ registry.Registrar = (*Registry)(nil)
	_ registry.Discovery = (*Registry)(nil)
)

// Registry is etcd registry.
type Registry struct {
	opts     *options
	client   *clientv3.Client
	kv       clientv3.KV
	lease    clientv3.Lease
	mu       sync.Mutex
	hbCancel context.CancelFunc
}

// New creates etcd registry
func New(client *clientv3.Client, opts ...Option) (r *Registry) {
	op := &options{
		ctx:       context.Background(),
		namespace: "/microservices",
		ttl:       time.Second * 15,
		maxRetry:  5,
	}
	for _, o := range opts {
		o(op)
	}
	return &Registry{
		opts:   op,
		client: client,
		kv:     clientv3.NewKV(client),
	}
}

// Register the registration.
func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	value, err := marshal(service)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)

	// create new lease locally first
	newLease := clientv3.NewLease(r.client)

	// swap leases under lock
	r.mu.Lock()
	oldLease := r.lease
	oldCancel := r.hbCancel
	r.lease = newLease
	r.hbCancel = nil
	r.mu.Unlock()

	// cancel old heartbeat and close old lease after swap
	if oldCancel != nil {
		oldCancel()
	}
	if oldLease != nil {
		_ = oldLease.Close()
	}

	// try to register with KV using ctx
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		// rollback: restore old lease and close the new one
		r.mu.Lock()
		r.lease = oldLease
		r.mu.Unlock()

		_ = newLease.Close()
		return err
	}

	// start heartbeat with registry-level context (long-living) and save cancel
	hbCtx, hbCancel := context.WithCancel(r.opts.ctx)
	r.mu.Lock()
	r.hbCancel = hbCancel
	r.mu.Unlock()
	go r.heartBeat(hbCtx, leaseID, key, value)

	return nil
}

// Deregister the registration.
func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	// remove kv first
	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, service.Name, service.ID)
	_, err := r.client.Delete(ctx, key)
	if err != nil {
		return wrapConnError("delete key", key, err)
	}

	// stop heartbeat and close lease safely
	r.mu.Lock()
	cancel := r.hbCancel
	l := r.lease
	r.hbCancel = nil
	r.lease = nil
	r.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if l != nil {
		_ = l.Close()
	}
	return err
}

// GetService return the service instances in memory according to the service name.
func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)

	resp, err := r.kv.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, wrapConnError("get key prefix", key, err)
	}

	items := make([]*registry.ServiceInstance, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		si, err := unmarshal(kv.Value)
		if err != nil {
			return nil, err
		}
		if si.Name != name {
			continue
		}
		items = append(items, si)
	}
	return items, nil
}

// Watch creates a watcher according to the service name.
func (r *Registry) Watch(ctx context.Context, name string) (registry.Watcher, error) {
	key := fmt.Sprintf("%s/%s", r.opts.namespace, name)
	w, err := newWatcher(ctx, key, name, r.client)
	if err != nil {
		return nil, wrapConnError("create watcher", key, err)
	}
	return w, nil
}

// registerWithKV create a new lease, return current leaseID
func (r *Registry) registerWithKV(ctx context.Context, key string, value string) (clientv3.LeaseID, error) {
	r.mu.Lock()
	l := r.lease
	r.mu.Unlock()

	if l == nil {
		return 0, errNoLease
	}

	grant, err := l.Grant(ctx, int64(r.opts.ttl.Seconds()))
	if err != nil {
		return 0, wrapConnError("grant lease", "", err)
	}

	_, err = r.kv.Put(ctx, key, value, clientv3.WithLease(grant.ID))
	if err != nil {
		return 0, wrapConnError("put key", key, err)
	}

	return grant.ID, nil
}

func (r *Registry) heartBeat(ctx context.Context, leaseID clientv3.LeaseID, key string, value string) {
	curLeaseID := leaseID
	kac, err := r.client.KeepAlive(ctx, leaseID)
	if err != nil {
		curLeaseID = 0
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		if curLeaseID == 0 {
			// try to registerWithKV
			var retreat []int
			for retryCnt := 0; retryCnt < r.opts.maxRetry; retryCnt++ {
				if ctx.Err() != nil {
					return
				}
				// prevent infinite blocking
				idChan := make(chan clientv3.LeaseID, 1)
				errChan := make(chan error, 1)
				cancelCtx, cancel := context.WithCancel(ctx)
				go func() {
					defer cancel()
					id, registerErr := r.registerWithKV(cancelCtx, key, value)
					if registerErr != nil {
						errChan <- registerErr
					} else {
						idChan <- id
					}
				}()

				select {
				case <-time.After(3 * time.Second):
					cancel()
					continue
				case <-errChan:
					continue
				case curLeaseID = <-idChan:
				}

				kac, err = r.client.KeepAlive(ctx, curLeaseID)
				if err == nil {
					break
				}
				retreat = append(retreat, 1<<retryCnt)
				time.Sleep(time.Duration(retreat[rnd.Intn(len(retreat))]) * time.Second)
			}
			if _, ok := <-kac; !ok {
				// retry failed
				return
			}
		}

		select {
		case _, ok := <-kac:
			if !ok {
				if ctx.Err() != nil {
					// channel closed due to context cancel
					return
				}
				// need to retry registration
				curLeaseID = 0
				continue
			}
		case <-r.opts.ctx.Done():
			return
		}
	}
}

// isConnError returns true when the error looks like a connection/server unreachable error.
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
		"tls:", "connection refused", "connection reset by peer", "eof",
	}
	for _, sub := range indicators {
		if strings.Contains(e, sub) {
			return true
		}
	}
	return false
}

// wrapConnError returns a clearer error message for connection related errors.
func wrapConnError(op string, key string, err error) error {
	if err == nil {
		return nil
	}
	if isConnError(err) {
		if key == "" {
			return fmt.Errorf("etcd: %s failed (cannot reach etcd server): %w", op, err)
		}
		return fmt.Errorf("etcd: %s failed for key %s (cannot reach etcd server): %w", op, key, err)
	}
	return fmt.Errorf("etcd: %s failed: %w", op, err)
}
