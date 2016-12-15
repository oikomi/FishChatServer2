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
	configFile             string
	Server                 *commconf.Server
	ServiceDiscoveryServer *commconf.ServiceDiscoveryServer
	ServiceDiscoveryClient *commconf.ServiceDiscoveryClient
	RPCServer              *commconf.RPCServer
	RPCClient              *RPCClient
	Etcd                   *commconf.Etcd
}

type RPCClient struct {
	LogicClient *commconf.RPCClient
}

func init() {
	flag.StringVar(&confPath, "conf", "./access.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
