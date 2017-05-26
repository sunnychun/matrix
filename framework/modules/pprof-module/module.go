package pprof_module

import (
	"context"
	"net"
	"net/http"

	_ "net/http/pprof"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/tlog"
)

var Config = &config{
	Addr: ":6060",
}

var Module = &module{}

func init() {
	framework.Register(Module, Config)
}

type config struct {
	Addr string
}

type module struct {
	ln net.Listener
}

func (m *module) Name() string {
	return "pprof-module"
}

func (m *module) Init() (err error) {
	log := tlog.Std().Sugar().With("module", m.Name())
	m.ln, err = net.Listen("tcp", Config.Addr)
	if err != nil {
		log.Errorw("listen", "error", err)
		return err
	}
	log.Debug("init success")
	return nil
}

func (m *module) Fini() error {
	return nil
}

func (m *module) Run(ctx context.Context) {
	log := tlog.Std().Sugar().With("module", m.Name(), "addr", Config.Addr)
	log.Debug("start serve")

	go func() {
		<-ctx.Done()
		m.ln.Close()
	}()
	http.Serve(m.ln, nil)

	log.Debug("stop serve")
}
