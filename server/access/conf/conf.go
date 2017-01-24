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
	RPCServer              *commconf.RPCServer
	RPCClient              *RPCClient
	ConfDiscovery          *ConfDiscovery
	// Etcd                   *commconf.Etcd
}

type RPCClient struct {
	LogicClient *commconf.ServiceDiscoveryClient
}

type ConfDiscovery struct {
	Gateway *commconf.Etcd
	MsgJob  *commconf.Etcd
	Notify  *commconf.Etcd
	Logic   *commconf.Etcd
}

func init() {
	flag.StringVar(&confPath, "conf", "./access.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
