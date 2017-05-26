package etcd_module

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/tlog"
)

var Config = &clientv3.Config{
	Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
	DialTimeout: 5 * time.Second,
}

var Module = &module{}

func init() {
	framework.Register(Module, Config)
}

type module struct {
	client *clientv3.Client
}

func (m *module) Name() string {
	return "etcd-module"
}

func (m *module) Init() (err error) {
	log := tlog.Std().Sugar().With("module", m.Name())
	m.client, err = clientv3.New(*Config)
	if err != nil {
		log.Errorw("new clientv3", "error", err)
		return err
	}
	log.Debug("init success")
	return nil
}

func (m *module) Fini() error {
	return m.client.Close()
}

func (m *module) Client() *clientv3.Client {
	return m.client
}
