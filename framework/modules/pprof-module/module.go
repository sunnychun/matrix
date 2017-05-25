package pprof_module

import (
	"context"
	"net"
	"net/http"

	_ "net/http/pprof"

	"github.com/ironzhang/matrix/framework"
)

var Module framework.Module = &module{}

type module struct {
	ln net.Listener
}

func (m *module) Name() string {
	return "pprof-module"
}

func (m *module) Init() (err error) {
	m.ln, err = net.Listen("tcp", ":6060")
	if err != nil {
		return err
	}
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
