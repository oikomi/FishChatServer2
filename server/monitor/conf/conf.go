package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	commconf "github.com/oikomi/FishChatServer2/common/conf"
)

var (
	confPath string
	Conf     *Config
)

type Config struct {
	*commconf.CommConf
	configFile string
	Server     *commconf.Server
	RPCServer  *commconf.RPCServer
	Etcd       *commconf.Etcd
	Zookeeper  *commconf.Zookeeper
}

func init() {
	flag.StringVar(&confPath, "conf", "./monitor.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
