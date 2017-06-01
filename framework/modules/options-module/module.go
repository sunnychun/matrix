package options_module

import (
	"flag"

	"github.com/ironzhang/matrix/framework/modules/options-module/meta"
	"github.com/ironzhang/matrix/framework/modules/options-module/options"
)

var CommandLine = flag.CommandLine

type M struct {
	values meta.Values
}

func (m *M) Name() string {
	return "options-module"
}

func (m *M) Init() error {
	return options.Setup(m.values)
	return nil
}

func (m *M) Fini() error {
	return nil
}

func (m *M) Register(name string, value interface{}) {
	if err := m.values.Register(name, value); err != nil {
		panic(err)
	}
}
