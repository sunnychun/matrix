package main

import (
	"bufio"
	"context"
	"encoding/json"
	"expvar"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/ironzhang/gomqtt/pkg/packet"
	"github.com/ironzhang/matrix/tlog"
)

type stats struct {
	Conns  int64
	Slows  int64
	Errors int64
}

func (s *stats) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

var g = stats{}

func init() {
	expvar.Publish("stats", &g)
}

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
			atomic.AddInt64(&g.Conns, 1)
			defer atomic.AddInt64(&g.Conns, -1)
			handleConn(ctx, c)
		}(c)
	}
}

func handleConn(ctx context.Context, c net.Conn) {
	go func() {
		<-ctx.Done()
		c.Close()
	}()

	//log := tlog.Std().Sugar().With("addr", c.RemoteAddr().String())

	r := bufio.NewReaderSize(c, 256)
	for {
		//c.SetReadDeadline(time.Now().Add(120 * time.Second))
		cp, err := packets.ReadPacket(r)
		if err != nil {
			if err != io.EOF {
				atomic.AddInt64(&g.Errors, 1)
			}
			break
		}

		start := time.Now()
		if err = processPacket(c, cp); err != nil {
			atomic.AddInt64(&g.Errors, 1)
			break
		}
		if time.Since(start) > time.Second {
			atomic.AddInt64(&g.Slows, 1)
		}
	}
}

func processPacket(c net.Conn, cp packets.ControlPacket) error {
	switch cp.(type) {
	case *packets.ConnectPacket:
		resp := packet.NewConnackPacket()
		resp.ReturnCode = packets.Accepted
		c.SetWriteDeadline(time.Now().Add(10 * time.Second))
		return resp.Write(c)
	case *packets.PingreqPacket:
		resp := packet.NewPingrespPacket()
		c.SetWriteDeadline(time.Now().Add(10 * time.Second))
		return resp.Write(c)
	}
	return nil
}
