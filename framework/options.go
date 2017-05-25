package framework

import "flag"

type Options struct {
	ConfigFile    string
	LogConfigFile string
}

func (o *Options) setup(f *flag.FlagSet) {
	f.StringVar(&o.ConfigFile, "-config-file", "./cfg.json", "config file")
	f.StringVar(&o.LogConfigFile, "-log-config-file", "", "log config file")
}
