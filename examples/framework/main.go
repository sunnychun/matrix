package main

import (
	"github.com/ironzhang/matrix/framework"

	_ "github.com/ironzhang/matrix/framework/modules/debug-module"
	_ "github.com/ironzhang/matrix/framework/modules/etcd-module"
	_ "github.com/ironzhang/matrix/framework/modules/micro-module"
)

var Module = &M{}

func init() {
	framework.Register(Module, nil)
}

type M struct {
}

func (m *M) Name() string {
	return "main-module"
}

func (m *M) Init() error {
	return nil
}

func (m *M) Fini() error {
	return nil
}

func main() {
	framework.Main()
}
