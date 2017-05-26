package pprof_module

import (
	"context"
	"net"
	"net/http"

	_ "net/http/pprof"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/tlog"
)

var Module framework.Module = &module{}

type module struct {
	ln net.Listener
}

func (m *module) Name() string {
	return "pprof-module"
}

func (m *module) Init() (err error) {
	log := tlog.Std().Sugar().With("module", m.Name())
	m.ln, err = net.Listen("tcp", Conf.Addr)
	if err != nil {
		return err
	}
	log.Debugw("listen", "addr", Conf.Addr)
	return nil
}

func (m *module) Fini() error {
	return nil
}

func (m *module) Serve(ctx context.Context) {
	go func() {
		<-ctx.Done()
		m.ln.Close()
	}()

	http.Serve(m.ln, nil)
}
