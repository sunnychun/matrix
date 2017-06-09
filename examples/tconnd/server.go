package main

import (
	"bufio"
	"context"
	"expvar"
	"fmt"
	"net"
	"time"

	"github.com/ironzhang/matrix/tlog"
)

var conns = expvar.NewInt("conns")
var slows = expvar.NewInt("slows")

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
	//log.Info("serve")
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

	r := bufio.NewReader(c)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		start := time.Now()
		if _, err = fmt.Fprintf(c, "%s\n", line); err != nil {
			break
		}
		if time.Since(start) > time.Second {
			slows.Add(1)
		}
	}
	c.Close()
}
