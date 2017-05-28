package debug_module

import (
	"net/http"
	"net/http/pprof"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/tlog"
)

var Module = &M{}

func init() {
	framework.Register(Module, nil)
}

type M struct {
}

func (m *M) Name() string {
	return "debug-module"
}

func (m *M) Init() (err error) {
	mux := http.NewServeMux()
	mux.Handle("/log/level", tlog.Level())
	mux.HandleFunc("/pprof/", pprof.Index)
	mux.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/pprof/profile", pprof.Profile)
	mux.HandleFunc("/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/pprof/trace", pprof.Trace)
	//backend_module.Handle("/debug", mux)
	return nil
}

func (m *M) Fini() error {
	return nil
}
