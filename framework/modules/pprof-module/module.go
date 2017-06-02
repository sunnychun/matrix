package pprof_module

import (
	"net/http"
	"net/http/pprof"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/backend-module"
	"github.com/ironzhang/matrix/tlog"
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
	mux := http.NewServeMux()
	mux.Handle("/debug/log/level", tlog.Level())
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug//pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	backend_module.Module.Handle("/debug/", mux)
	return nil
}

func (m *M) Fini() error {
	return nil
}
