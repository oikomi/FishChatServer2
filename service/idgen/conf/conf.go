package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	commconf "github.com/oikomi/FishChatServer2/common/conf"
)

var (
	confPath string
	Conf     *Config
	Etcd     *commconf.Etcd
)

type Config struct {
	*commconf.CommConf
	configFile             string
	RPCServer              *commconf.RPCServer
	ServiceDiscoveryServer *commconf.ServiceDiscoveryServer
	Etcd                   *commconf.Etcd
}

func init() {
	flag.StringVar(&confPath, "conf", "./idgen.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
