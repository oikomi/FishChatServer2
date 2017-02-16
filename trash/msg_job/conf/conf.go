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
	MongoDB       *MongoDB
	Etcd          *commconf.Etcd
}

type RPCClient struct {
	AccessClient *commconf.ServiceDiscoveryClient
}

type MongoDB struct {
	*commconf.MongoDB
	OfflineMsgCollection string
}

func init() {
	flag.StringVar(&confPath, "conf", "./msg_job.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
