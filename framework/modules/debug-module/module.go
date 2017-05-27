package debug_module

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/tlog"
)

var Config = &C{
	Addr: ":6060",
}

var Module = &M{}

func init() {
	framework.Register(Module, Config)
}

type C struct {
	Addr string
}

type M struct {
	ln net.Listener
}

func (m *M) Name() string {
	return "debug-module"
}

func (m *M) Init() (err error) {
	log := tlog.Std().Sugar().With("module", m.Name())
	m.ln, err = net.Listen("tcp", Config.Addr)
	if err != nil {
		log.Errorw("listen", "error", err)
		return err
	}
	log.Debug("init success")
	return nil
}

func (m *M) Fini() error {
	return nil
}

func (m *M) Run(ctx context.Context) {
	log := tlog.Std().Sugar().With("module", m.Name(), "addr", Config.Addr)
	log.Debug("start serve")

	go func() {
		<-ctx.Done()
		m.ln.Close()
	}()

	mux := http.NewServeMux()
	mux.Handle("/debug/log/level", tlog.Level())
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	http.Serve(m.ln, mux)

	log.Debug("stop serve")
}
