package main

import (
	"bufio"
	"context"
	"expvar"
	"fmt"
	"net"

	"github.com/ironzhang/matrix/tlog"
)

var conns = expvar.NewInt("conns")

type server struct {
	ln net.Listener
}

func (s *server) Init(addr string) (err error) {
	s.ln, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) serve(ctx context.Context) {
	go func() {
		<-ctx.Done()
		s.ln.Close()
	}()

	log := tlog.Std().Sugar().With("addr", s.ln.Addr())
	log.Info("serve")
	for {
		c, err := s.ln.Accept()
		if err != nil {
			log.Infow("accpet", "error", err)
			return
		}

		go func(c net.Conn) {
			conns.Add(1)
			defer conns.Add(-1)
			handleConn(ctx, c)
		}(c)
	}
}

func handleConn(ctx context.Context, c net.Conn) {
	go func() {
		<-ctx.Done()
		c.Close()
	}()

	input := bufio.NewScanner(c)
	for input.Scan() {
		fmt.Fprintln(c, input.Text())
	}
	c.Close()
}
