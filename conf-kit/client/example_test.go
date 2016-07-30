package client_test

import "github.com/ironzhang/matrix/conf-kit/client"

type EnvConfig struct {
	EtcdAddrs []string
}

type ServiceConfig struct {
	MysqlHost string
	MysqlPort string
}

type ProcessConfig struct {
	Host string
	Port string
}

type Config struct {
	Env     EnvConfig
	Service ServiceConfig
	Process ProcessConfig
}

func ExampleClient() {
	addrs := []string{"127.0.0.1:7200", "127.0.0.1:7201"}
	c := client.New(addrs)

	var conf Config
	c.LoadConfig("/dev", &conf.Env, nil)
	c.LoadConfig("/dev/ac-account", &conf.Service, ReloadServiceConfig)
	c.LoadConfig("/dev/ac-account/p1", &conf.Process, nil)

	c.LoadFile("local_env_file.json", "/dev", &conf.Env, nil)
	c.LoadFile("local_service_file.json", "/dev/ac-account", &conf.Service, ReloadServiceConfig)
	c.LoadFile("local_process_file.json", "/dev/ac-account/p1", &conf.Process, nil)
}

func ReloadServiceConfig(cfg *ServiceConfig) {
}
