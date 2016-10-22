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
	RPCClient     *RPCClient
	KafkaConsumer *commconf.KafkaConsumer
	Etcd          *commconf.Etcd
}

type RPCClient struct {
	MsgServerClient *commconf.RPCClient
}

func init() {
	flag.StringVar(&confPath, "conf", "./msg-job.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
