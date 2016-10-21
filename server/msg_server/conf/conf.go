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
	RPCServer     *commconf.RPCServer
	RPCClient     *RPCClient
	Etcd          *commconf.Etcd
	KafkaProducer *KafkaProducer
	Zookeeper     *commconf.Zookeeper
}

type KafkaProducer struct {
	Topic    string
	Producer *commconf.KafkaProducer
}

type RPCClient struct {
	ManagerClient *commconf.RPCClient
}

func init() {
	flag.StringVar(&confPath, "conf", "./msg_server.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
