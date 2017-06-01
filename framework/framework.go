package framework

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ironzhang/matrix/framework/config"
	"github.com/ironzhang/matrix/tlog"
)

type Module interface {
	Name() string
	Init() error
	Fini() error
}

type Runner interface {
	Run(ctx context.Context)
}

type framework struct {
	options options
	configs config.ConfigSet
	modules []Module
}

func (f *framework) doCommandLine() {
	var err error
	var quit bool

	if f.options.ConfigExample != "" {
		if err = f.configs.WriteToFile(f.options.ConfigExample); err != nil {
			fmt.Fprintf(os.Stderr, "generate config example: %v\n", err)
			os.Exit(3)
		}
		quit = true
	}
	if f.options.LogConfigExample != "" {
		if err = tlogWriteToFile(f.options.LogConfigExample); err != nil {
			fmt.Fprintf(os.Stderr, "generate log config example: %v\n", err)
			os.Exit(3)
		}
		quit = true
	}

	if quit {
		os.Exit(0)
	}
}

func (f *framework) main() (err error) {
	log := tlog.Std().Sugar()

	// module init
	for _, m := range f.modules {
		if err = m.Init(); err != nil {
			log.Errorw("init", "module", m.Name(), "error", err)
			return fmt.Errorf("init %s: %v", m.Name(), err)
		}
		log.Debugw("init", "module", m.Name())
	}

	// quit signal
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
		cancel()
		time.Sleep(10 * time.Second)
		fmt.Fprintf(os.Stderr, "wait 10s, force exit")
		log.Sync()
		os.Exit(-3)
	}()

	// module run
	var wg sync.WaitGroup
	for _, m := range f.modules {
		if r, ok := m.(Runner); ok {
			wg.Add(1)
			go func(r Runner) {
				defer wg.Done()
				r.Run(ctx)
			}(r)
		}
	}
	wg.Wait()

	// module fini
	for i := len(f.modules) - 1; i >= 0; i-- {
		m := f.modules[i]
		if err = m.Fini(); err != nil {
			log.Errorw("fini", "module", m.Name(), "error", err)
			continue
		}
		log.Debugw("fini", "module", m.Name())
	}
	return nil
}

func (f *framework) Main() {
	var err error
	f.options.Parse()
	f.doCommandLine()

	// load config
	if f.options.ConfigFile != "" {
		if err = f.configs.LoadFromFile(f.options.ConfigFile); err != nil {
			fmt.Fprintf(os.Stderr, "load config: %v\n", err)
			os.Exit(3)
		}
	}

	// tlog load from file
	log, err := tlogLoadFromFile(f.options.LogConfigFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tlog load from file: %v\n", err)
		os.Exit(3)
	}
	defer log.Sync()

	// main
	if err = f.main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}

func (f *framework) Register(m Module, config interface{}) {
	for _, v := range f.modules {
		if v.Name() == m.Name() {
			panic(fmt.Sprintf("module(%s) duplicate", m.Name()))
		}
	}
	f.modules = append(f.modules, m)

	if config != nil {
		if err := f.configs.Register(m.Name(), config); err != nil {
			panic(err)
		}
	}
}

func (f *framework) ModuleConfig(m string) (interface{}, bool) {
	return f.configs.Config(m)
}

var f = &framework{}

func Main() {
	f.Main()
}

func Register(m Module, config interface{}) {
	f.Register(m, config)
}

func ModuleConfig(m string) (interface{}, bool) {
	return f.ModuleConfig(m)
}
