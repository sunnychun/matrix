package main

import (
	"github.com/ironzhang/matrix/framework"

	_ "github.com/ironzhang/matrix/framework/modules/dashboard-module"
	_ "github.com/ironzhang/matrix/framework/modules/pprof-module"
)

var Options = &O{}

var Module = &M{}

func init() {
	framework.Register(Module, Options, nil)
}

type O struct {
	Addr string `json:",readonly" usage:"指定监听地址"`
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
