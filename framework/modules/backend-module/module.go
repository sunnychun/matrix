package backend_module

import (
	"context"
	"net"
	"net/http"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/tlog"
)

var Config = &C{
	Addr: ":6060",
}

var Module = &M{}

func init() {
	framework.Register(Module, nil, Config)
}

type C struct {
	Addr string
}

type M struct {
	http.ServeMux
	ln net.Listener
}

func (m *M) Name() string {
	return "backend-module"
}

func (m *M) Init() (err error) {
	if m.ln, err = net.Listen("tcp", Config.Addr); err != nil {
		return err
	}
	return nil
}

func (m *M) Fini() error {
	return nil
}

func (m *M) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		m.ln.Close()
	}()

	log := tlog.Std().Sugar().With("module", m.Name(), "addr", Config.Addr)
	log.Info("start")
	http.Serve(m.ln, &m.ServeMux)
	log.Info("stop")
}
