package micro_module

import (
	"time"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/etcd-module"
	"github.com/ironzhang/matrix/micro/discovery"
	"github.com/ironzhang/matrix/micro/registry"
	"github.com/ironzhang/matrix/tlog"
)

var Config = &config{
	Namespace: "matrix",
	TTL:       10,
	Timeout:   5 * time.Second,
}

var Module = &module{}

func init() {
	framework.Register(Module, Config)
}

type config struct {
	Namespace string
	Timeout   time.Duration
	TTL       int64
}

type module struct {
	r *registry.Registry
	d *discovery.Discovery
}

func (m *module) Name() string {
	return "micro-module"
}

func (m *module) Init() error {
	log := tlog.Std().Sugar().With("module", m.Name())
	c := etcd_module.Module.Client()
	m.r = registry.New(c, registry.Options{TTL: Config.TTL, Timeout: Config.Timeout, Namespace: Config.Namespace})
	m.d = discovery.New(c, discovery.Options{Namespace: Config.Namespace, Timeout: Config.Timeout})
	log.Debug("init success")
	return nil
}

func (m *module) Fini() error {
	m.r.UnregisterAll()
	m.d.UnwatchAll()
	return nil
}

func (m *module) Registry() *registry.Registry {
	return m.r
}

func (m *module) Discovery() *discovery.Discovery {
	return m.d
}
