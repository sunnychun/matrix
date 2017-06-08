package pprof_module

import (
	"net/http/pprof"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/backend-module"
)

var Module = &M{}

func init() {
	framework.Register(Module, nil, nil)
}

type M struct {
}

func (m *M) Name() string {
	return "pprof-module"
}

func (m *M) Init() error {
	backend_module.Module.HandleFunc("/debug/pprof/", pprof.Index)
	backend_module.Module.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	backend_module.Module.HandleFunc("/debug//pprof/profile", pprof.Profile)
	backend_module.Module.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	backend_module.Module.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return nil
}

func (m *M) Fini() error {
	return nil
}
