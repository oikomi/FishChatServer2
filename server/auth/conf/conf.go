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
	Auth       *Auth
	RPCServer  *commconf.RPCServer
	Etcd       *commconf.Etcd
}

type Auth struct {
	Encryption string
	Salt       string
}

func init() {
	flag.StringVar(&confPath, "conf", "./auth.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
