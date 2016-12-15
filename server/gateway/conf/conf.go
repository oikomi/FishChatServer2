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
	configFile    string
	Server        *commconf.Server
	ConfDiscovery *commconf.ConfDiscovery
	Etcd          *commconf.Etcd
}

func init() {
	flag.StringVar(&confPath, "conf", "./gateway.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
