package framework

import "flag"

type Options struct {
	ConfigFile    string
	LogConfigFile string
	ConfigExample string
}

func (o *Options) setup(f *flag.FlagSet) {
	f.StringVar(&o.ConfigFile, "config-file", "", "指定配置文件")
	f.StringVar(&o.LogConfigFile, "log-config-file", "", "指定日志配置文件")
	f.StringVar(&o.ConfigExample, "config-example", "", "指定生成的示例配置文件")
}
