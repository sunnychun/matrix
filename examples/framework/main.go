package main

import (
	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/tlog"

	_ "github.com/ironzhang/matrix/framework/modules/pprof-module"
)

type module struct {
}

func (m *module) Name() string {
	return "main-module"
}

func (m *module) Init() error {
	log := tlog.Std().Sugar().With("module", m.Name())
	log.Debug("init")
	return nil
}

func (m *module) Fini() error {
	log := tlog.Std().Sugar().With("module", m.Name())
	log.Debug("fini")
	return nil
}

var Module = &module{}

func init() {
	framework.Register(Module, nil)
}

func main() {
	framework.Main()
}
