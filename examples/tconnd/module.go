package main

import (
	"context"
	"sync"

	"github.com/ironzhang/matrix/framework"
	"github.com/ironzhang/matrix/tlog"
)

var Config = &C{
	Addrs: []string{":1883", ":2000", ":2001", ":2002", ":2003", ":2004", ":2005", ":2006", ":2007", ":2008"},
}

var Module = &M{}

func init() {
	framework.Register(Module, nil, Config)
}

type C struct {
	Addrs []string `json:",readonly"`
}

type M struct {
	servers []server
}

func (m *M) Name() string {
	return "main-module"
}

func (m *M) Init() (err error) {
	m.servers = make([]server, len(Config.Addrs))
	for i, addr := range Config.Addrs {
		if err = m.servers[i].Init(addr); err != nil {
			return err
		}
	}
	return nil
}

func (m *M) Fini() error {
	return nil
}

func (m *M) Run(ctx context.Context) {
	log := tlog.Std().Sugar().With("module", m.Name())

	log.Info("start")
	var wg sync.WaitGroup
	for i := range m.servers {
		wg.Add(1)
		go func(s *server) {
			defer wg.Done()
			s.serve(ctx)
		}(&m.servers[i])
	}
	wg.Wait()
	log.Info("stop")
}
