package experimental

import (
	"context"
	"testing"
	"time"

	"github.com/coreos/etcd/cmd/etcd/clientv3"
)

func get(c *clientv3.Client, key string) (string, error) {
	resp, err := c.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), nil
	}
	return "", nil
}

func KVAssert(t *testing.T, n string, c *clientv3.Client, key, value string) {
	got, err := get(c, key)
	if err != nil {
		t.Fatalf("%s: %v", n, err)
	}
	if want := value; got != want {
		t.Errorf("%s: key: %q, value: %q != %q", n, key, got, want)
	}
}

func TestRegister(t *testing.T) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	key := "sample_key"
	value := "sample_value"
	ttl := int64(2)

	id, err := register(context.Background(), c, key, value, ttl)
	if err != nil {
		t.Fatal(err)
	}
	c.KeepAlive(context.Background(), id)
	KVAssert(t, "register", c, key, value)
	time.Sleep(time.Duration(ttl) * time.Second)
	KVAssert(t, "sleep", c, key, value)
	time.Sleep(time.Duration(ttl) * time.Second)
	KVAssert(t, "sleep sleep", c, key, value)
	time.Sleep(time.Duration(ttl) * time.Second)
	KVAssert(t, "sleep sleep sleep", c, key, value)

	unregister(context.Background(), c, id, key)
	KVAssert(t, "unregister", c, key, "")
}
