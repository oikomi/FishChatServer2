package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	commconf "github.com/oikomi/FishChatServer2/common/conf"
	"time"
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
	RPCClient  *RPCClient
	Etcd       *commconf.Etcd
	Zookeeper  *commconf.Zookeeper
	Redis      *Redis
}

type RPCClient struct {
	MsgServerClient *commconf.RPCClient
}

type Redis struct {
	*commconf.Redis
	Expire time.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "./manager.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
