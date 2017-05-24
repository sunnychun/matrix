package registry_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/micro/registry"
	"github.com/ironzhang/matrix/tlog"
)

func NewClient(t *testing.T) *clientv3.Client {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func endpoint(namespace string, p registry.Endpoint) string {
	return fmt.Sprintf("%s/%s/%s", namespace, p.Service, p.Addr)
}

func KeyIsExist(c *clientv3.Client, key string) bool {
	log := tlog.Std().Sugar().With("key", key)
	resp, err := c.Get(context.Background(), key)
	if err != nil {
		log.Errorw("client get", "error", err)
		return false
	}
	if len(resp.Kvs) != 1 {
		return false
	}
	if string(resp.Kvs[0].Key) != key {
		return false
	}
	return true
}

func TestRegistry(t *testing.T) {
	c := NewClient(t)

	ns := "TestRegistry"
	p1 := registry.Endpoint{"S1", "A1"}
	p2 := registry.Endpoint{"S2", "A2"}

	r, err := registry.New(c, registry.Options{Namespace: ns})
	if err != nil {
		t.Fatalf("new: %v", err)
	}

	if key := endpoint(ns, p1); KeyIsExist(c, key) {
		t.Fatalf("before register: key(%s) is existed", key)
	}
	if key := endpoint(ns, p2); KeyIsExist(c, key) {
		t.Fatalf("before register: key(%s) is existed", key)
	}

	if err = r.Register(p1); err != nil {
		t.Fatalf("register p1: %v", err)
	}
	if err = r.Register(p2); err != nil {
		t.Fatalf("register p2: %v", err)
	}
	if key := endpoint(ns, p1); !KeyIsExist(c, key) {
		t.Fatalf("after register: key(%s) is not existed", key)
	}
	if key := endpoint(ns, p2); !KeyIsExist(c, key) {
		t.Fatalf("after register: key(%s) is not existed", key)
	}

	if err = r.Unregister(p1); err != nil {
		t.Fatalf("unregister p1: %v", err)
	}
	if key := endpoint(ns, p1); KeyIsExist(c, key) {
		t.Fatalf("after unregister p1: key(%s) is existed", key)
	}
	if key := endpoint(ns, p2); !KeyIsExist(c, key) {
		t.Fatalf("after unregister p1: key(%s) is not existed", key)
	}

	if err = r.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if key := endpoint(ns, p1); KeyIsExist(c, key) {
		t.Fatalf("after close: key(%s) is existed", key)
	}
	if key := endpoint(ns, p2); KeyIsExist(c, key) {
		t.Fatalf("after close: key(%s) is existed", key)
	}
}
