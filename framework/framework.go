package framework

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

type Module interface {
	Name() string
	Init() error
	Fini() error
	Serve(ctx context.Context)
}

type Framework struct {
	FlagSet *flag.FlagSet

	LoadConfigFunc func(file string) error

	OnInitFunc func() error

	OnFiniFunc func() error

	Modules []Module

	options Options
}

func (f *Framework) init() {
	if f.FlagSet == nil {
		f.FlagSet = flag.CommandLine
	}
}

func (f *Framework) loadConfig(file string) error {
	if f.LoadConfigFunc != nil {
		return f.LoadConfigFunc(file)
	}
	return nil
}

func (f *Framework) onInit() error {
	if f.OnInitFunc != nil {
		return f.OnInitFunc()
	}
	return nil
}

func (f *Framework) onFini() error {
	if f.OnFiniFunc != nil {
		return f.OnFiniFunc()
	}
	return nil
}

func (f *Framework) Main() {
	f.init()
	f.options.setup(f.FlagSet)
	f.FlagSet.Parse(os.Args[1:])

	// tlog load from file
	log, err := tlogLoadFromFile(f.options.LogConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tlog load from file: %v\n", err)
		os.Exit(3)
	}
	defer log.Sync()

	// load config
	if err = f.loadConfig(f.options.ConfigFile); err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(3)
	}

	// on init
	if err = f.onInit(); err != nil {
		fmt.Fprintf(os.Stderr, "on init: %v\n", err)
		os.Exit(3)
	}

	// module init
	for _, m := range f.Modules {
		if err = m.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "module(%s) init: %v\n", m.Name(), err)
			os.Exit(3)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
		cancel()
	}()

	// module serve
	var wg sync.WaitGroup
	for _, m := range f.Modules {
		wg.Add(1)
		go func(m Module) {
			defer wg.Done()
			m.Serve(ctx)
		}(m)
	}
	wg.Wait()

	// module fini
	var code int
	for _, m := range f.Modules {
		if err = m.Fini(); err != nil {
			code = -3
			fmt.Fprintf(os.Stderr, "module(%s) fini: %v\n", m.Name(), err)
			continue
		}
	}

	// on fini
	if err = f.onFini(); err != nil {
		fmt.Fprintf(os.Stderr, "on fini: %v\n", err)
		os.Exit(3)
	}

	os.Exit(code)
}
