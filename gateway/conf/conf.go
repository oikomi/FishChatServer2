package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/oikomi/FishChatServer2/common"
)

var (
	confPath string
	Conf     *Config
)

type Config struct {
	*common.CommConf
	configFile string
	Server     *common.Server
	Listen     string
}

func init() {
	flag.StringVar(&confPath, "conf", "./gateway.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
