package pprof_module

import "github.com/ironzhang/matrix/framework"

type conf struct {
	Addr string
}

var Conf = conf{
	Addr: ":6060",
}

func init() {
	framework.Config.Register(Module.Name(), &Conf)
}
