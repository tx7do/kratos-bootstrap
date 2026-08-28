package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"google.golang.org/grpc"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/go-kratos/kratos/v2/registry"
)

func createTestEtcdClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
}

func TestRegistry(t *testing.T) {
	client, err := createTestEtcdClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	s := &registry.ServiceInstance{
		ID:   "0",
		Name: "helloworld",
	}

	r := New(client)
	w, err := r.Watch(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = w.Stop()
	}()
	go func() {
		for {
			res, err1 := w.Next()
			if err1 != nil {
				return
			}
			t.Logf("watch: %d", len(res))
			for _, r := range res {
				t.Logf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)

	if err1 := r.Register(ctx, s); err1 != nil {
		t.Fatal(err1)
	}
	time.Sleep(time.Second)

	res, err := r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 && res[0].Name != s.Name {
		t.Errorf("not expected: %+v", res)
	}

	if err1 := r.Deregister(ctx, s); err1 != nil {
		t.Fatal(err1)
	}
	time.Sleep(time.Second)

	res, err = r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Errorf("not expected empty")
	}
}

func TestHeartBeat(t *testing.T) {
	client, err := createTestEtcdClient()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	s := &registry.ServiceInstance{
		ID:   "0",
		Name: "helloworld",
	}

	go func() {
		r := New(client)
		w, err1 := r.Watch(ctx, s.Name)
		if err1 != nil {
			return
		}
		defer func() {
			_ = w.Stop()
		}()
		for {
			res, err2 := w.Next()
			if err2 != nil {
				return
			}
			t.Logf("watch: %d", len(res))
			for _, r := range res {
				t.Logf("next: %+v", r)
			}
		}
	}()
	time.Sleep(time.Second)

	// new a server
	r := New(client,
		RegisterTTL(2*time.Second),
		MaxRetry(5),
	)

	key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, s.Name, s.ID)
	value, _ := marshal(s)
	r.lease = clientv3.NewLease(r.client)
	leaseID, err := r.registerWithKV(ctx, key, value)
	if err != nil {
		t.Fatal(err)
	}

	// wait for lease expired
	time.Sleep(3 * time.Second)

	res, err := r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 0 {
		t.Errorf("not expected empty")
	}

	go r.heartBeat(ctx, leaseID, key, value)

	time.Sleep(time.Second)
	res, err = r.GetService(ctx, s.Name)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) == 0 {
		t.Errorf("reconnect failed")
	}
}

// TestHeartBeatReRegister covers the 2026-08 incident: the service keeps
// running while its registration disappears. Two flavors:
//   - LeaseRevoked: the lease is revoked under a healthy etcd — the heartbeat
//     must re-register promptly.
//   - OutageOutlastingBurst: grant attempts keep failing (etcd unreachable)
//     for longer than one full retry burst — the heartbeat must keep retrying
//     and re-register once etcd is reachable again. The pre-fix heartbeat
//     returned permanently here, leaving the running service de-registered
//     forever.
func TestHeartBeatReRegister(t *testing.T) {
	setup := func(t *testing.T, name string) (*Registry, clientv3.LeaseID) {
		t.Helper()
		client, err := createTestEtcdClient()
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = client.Close() })

		ctx := context.Background()
		s := &registry.ServiceInstance{ID: "0", Name: name}

		r := New(client,
			RegisterTTL(2*time.Second),
			MaxRetry(5),
		)

		key := fmt.Sprintf("%s/%s/%s", r.opts.namespace, s.Name, s.ID)
		value, _ := marshal(s)
		r.lease = clientv3.NewLease(r.client)
		leaseID, err := r.registerWithKV(ctx, key, value)
		if err != nil {
			t.Fatal(err)
		}

		hbCtx, hbCancel := context.WithCancel(ctx)
		go r.heartBeat(hbCtx, leaseID, key, value)
		t.Cleanup(hbCancel)

		waitPresent(t, r, ctx, name, true, 5*time.Second)
		return r, leaseID
	}

	t.Run("LeaseRevoked", func(t *testing.T) {
		r, leaseID := setup(t, "helloworld-rereg-revoke")

		if _, err := r.client.Revoke(context.Background(), leaseID); err != nil {
			t.Fatal(err)
		}
		waitPresent(t, r, context.Background(), "helloworld-rereg-revoke", true, 30*time.Second)
	})

	t.Run("OutageOutlastingBurst", func(t *testing.T) {
		ctx := context.Background()
		r, leaseID := setup(t, "helloworld-rereg-outage")

		old := heartBeatRetryCooldown
		heartBeatRetryCooldown = 2 * time.Second
		t.Cleanup(func() { heartBeatRetryCooldown = old })

		// point the registry's lease at an unreachable endpoint so Grant
		// hangs until the per-attempt timeout — "etcd is down"
		badClient, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"127.0.0.1:1"},
			DialTimeout: 500 * time.Millisecond,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = badClient.Close() })

		r.mu.Lock()
		goodLease := r.lease
		r.lease = clientv3.NewLease(badClient)
		r.mu.Unlock()

		// lose the live registration while etcd is "down"
		if _, err := r.client.Revoke(ctx, leaseID); err != nil {
			t.Fatal(err)
		}

		// let a full retry burst (maxRetry x 3s per-attempt timeout) exhaust
		time.Sleep(20 * time.Second)

		// etcd "comes back": restore the working lease
		r.mu.Lock()
		r.lease = goodLease
		r.mu.Unlock()

		waitPresent(t, r, ctx, "helloworld-rereg-outage", true, 90*time.Second)
	})
}

func waitPresent(t *testing.T, r *Registry, ctx context.Context, name string, want bool, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		res, err := r.GetService(ctx, name)
		if err == nil {
			found := false
			for _, it := range res {
				if it.ID == "0" {
					found = true
					break
				}
			}
			if found == want {
				return
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("expected instance present=%v within %s", want, timeout)
}
