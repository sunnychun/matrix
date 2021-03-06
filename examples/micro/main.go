package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/micro/discovery"
	"github.com/ironzhang/matrix/micro/registry"
	"github.com/ironzhang/matrix/tlog"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "localhost:2000", "address")
	flag.Parse()

	tlog.Init(tlog.Config{Level: zap.DebugLevel, DisableStacktrace: true})

	c, err := NewEtcdClient()
	if err != nil {
		fmt.Printf("new etcd client: %v\n", err)
		return
	}

	r := registry.New(c, registry.Options{Timeout: 2 * time.Second})
	defer r.UnregisterAll()

	d := discovery.New(c, discovery.Options{Timeout: 2 * time.Second})
	defer d.UnwatchAll()

	// register
	p := registry.Endpoint{"micro", addr}
	if err = r.Register(p); err != nil {
		fmt.Printf("register: %v\n", err)
		return
	}

	// watch
	if _, err = d.Watch("micro", Refresh); err != nil {
		fmt.Printf("watch: %v\n", err)
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func Refresh(addrs []string) {
	fmt.Printf("refresh:\n")
	for _, addr := range addrs {
		fmt.Printf("\t%s\n", addr)
	}
	fmt.Printf("\n")
}

func NewEtcdClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
}
