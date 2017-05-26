package main

import (
	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/pprof-module"
	"github.com/ironzhang/matrix/tlog"
)

func OnInit() error {
	log := tlog.Std()
	log.Debug("on init")
	return nil
}

func OnFini() error {
	log := tlog.Std()
	log.Debug("on fini")
	return nil
}

type testModule struct {
}

func (m *testModule) Name() string {
	return "TestModule"
}

func (m *testModule) Init() error {
	log := tlog.Std()
	log.Debug("test module init")
	return nil
}

func (m *testModule) Fini() error {
	log := tlog.Std()
	log.Debug("test module fini")
	return nil
}

var TestModule = &testModule{}

func main() {
	(&framework.Framework{
		OnInitFunc: OnInit,
		OnFiniFunc: OnFini,
		Modules:    []framework.Module{TestModule, pprof_module.Module},
	}).Main()
}
