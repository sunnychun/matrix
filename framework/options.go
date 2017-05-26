package framework

import (
	"flag"
	"os"
)

var CommandLine = flag.CommandLine

type options struct {
	ConfigFile       string
	ConfigExample    string
	LogConfigFile    string
	LogConfigExample string
}

func (o *options) Parse() {
	CommandLine.StringVar(&o.ConfigFile, "config-file", "", "指定配置文件选项")
	CommandLine.StringVar(&o.ConfigExample, "config-example", "", "生成配置示例选项")
	CommandLine.StringVar(&o.LogConfigFile, "log-config-file", "", "指定日志配置文件选项")
	CommandLine.StringVar(&o.LogConfigExample, "log-config-example", "", "生成日志配置示例选项")
	CommandLine.Parse(os.Args[1:])
}
