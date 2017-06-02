package etcd_module

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/framework"
)

var Config = &clientv3.Config{
	Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
	DialTimeout: 5 * time.Second,
}

var Module = &M{}

func init() {
	framework.Register(Module, nil, Config)
}

type M struct {
	client *clientv3.Client
}

func (m *M) Name() string {
	return "etcd-module"
}

func (m *M) Init() (err error) {
	m.client, err = clientv3.New(*Config)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	return nil
}

func (m *M) Fini() error {
	return m.client.Close()
}

func (m *M) Client() *clientv3.Client {
	return m.client
}
