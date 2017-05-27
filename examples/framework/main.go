package main

import (
	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/tlog"

	_ "github.com/ironzhang/matrix/framework/modules/debug-module"
	_ "github.com/ironzhang/matrix/framework/modules/etcd-module"
	_ "github.com/ironzhang/matrix/framework/modules/micro-module"
)

var Module = &module{}

func init() {
	framework.Register(Module, nil)
}

type module struct {
}

func (m *module) Name() string {
	return "main-module"
}

func (m *module) Init() error {
	log := tlog.Std().Sugar().With("module", m.Name())
	log.Debug("init success")
	return nil
}

func (m *module) Fini() error {
	log := tlog.Std().Sugar().With("module", m.Name())
	log.Debug("fini success")
	return nil
}

func main() {
	framework.Main()
}
