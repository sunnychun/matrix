package dashboard_module

import (
	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/framework/modules/backend-module"
	"github.com/ironzhang/matrix/restful"
)

var Module = &M{}

func init() {
	framework.Register(Module, nil, nil)
}

type M struct {
}

func (m *M) Name() string {
	return "dashboard-module"
}

func (m *M) Init() (err error) {
	var h handlers
	mux := restful.NewServeMux(nil)
	mux.SetVerbose(2)
	if err = h.Register(mux); err != nil {
		return err
	}
	backend_module.Module.Handle("/config/", mux)
	return nil
}

func (m *M) Fini() error {
	return nil
}
