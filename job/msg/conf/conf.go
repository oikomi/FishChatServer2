package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/oikomi/FishChatServer2/common/conf"
	commconf "github.com/oikomi/FishChatServer2/common/conf"
)

var (
	confPath string
	Conf     *Config
)

type Config struct {
	*commconf.CommConf
	KafkaProducer *KafkaProducer
	KafkaConsumer *conf.KafkaConsumer
}

type KafkaProducer struct {
	Topic    string
	Producer *conf.KafkaProducer
}

func init() {
	flag.StringVar(&confPath, "conf", "./msg.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
