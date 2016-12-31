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
	Redis                  *Redis
	MongoDB                *MongoDB
}

type Redis struct {
	*commconf.Redis
	Expire xtime.Duration
}

type MongoDB struct {
	*commconf.MongoDB
	OfflineMsgCollection string
}

func init() {
	flag.StringVar(&confPath, "conf", "./manager.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
