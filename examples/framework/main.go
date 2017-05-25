package main

import (
	"context"
	"flag"
	"os"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/pprof-module"
	"github.com/ironzhang/matrix/tlog"
)

func LoadConfig(file string) error {
	log := tlog.Std().Sugar()
	log.Debugw("load config", "flie", file)
	return nil
}

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

func (m *testModule) Serve(ctx context.Context) {
	log := tlog.Std()
	log.Debug("test module serve start")
	<-ctx.Done()
	log.Debug("test module serve stop")
}

var TestModule = &testModule{}

func main() {
	(&framework.Framework{
		FlagSet:        flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		LoadConfigFunc: LoadConfig,
		OnInitFunc:     OnInit,
		OnFiniFunc:     OnFini,
		Modules:        []framework.Module{TestModule, pprof_module.Module},
	}).Main()
}
