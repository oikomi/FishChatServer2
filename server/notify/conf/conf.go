package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	commconf "github.com/oikomi/FishChatServer2/common/conf"
	"github.com/oikomi/FishChatServer2/common/xtime"
)

var (
	confPath string
	Conf     *Config
)

type Config struct {
	*commconf.CommConf
	configFile             string
	RPCServer              *commconf.RPCServer
	ServiceDiscoveryServer *commconf.ServiceDiscoveryServer
	RPCClient              *RPCClient
	Redis                  *Redis
	Mysql                  *Mysql
	Etcd                   *commconf.Etcd
}

type RPCClient struct {
	AccessClient *commconf.ServiceDiscoveryClient
}

type Redis struct {
	*commconf.Redis
	Expire xtime.Duration
}

type Mysql struct {
	IM *commconf.MySQL
}

// type MongoDB struct {
// 	*commconf.MongoDB
// 	GroupCollection string
// }

func init() {
	flag.StringVar(&confPath, "conf", "./notify.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
