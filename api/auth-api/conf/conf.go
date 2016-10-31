package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/oikomi/FishChatServer2/common/conf"
	"github.com/oikomi/FishChatServer2/common/xtime"
)

var (
	confPath string
	Conf     *Config
)

type Config struct {
	conf.CommConf
	MultiHTTP *conf.MultiHTTP
	Redis     *Redis
}

type Redis struct {
	*conf.Redis
	Expire xtime.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "./auth-api.toml", "config path")
}

func Init() error {
	if _, err := toml.DecodeFile(confPath, &Conf); err != nil {
		return err
	}
	return nil
}
