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
	MongoDB       *MongoDB
	// ES            *ES
}

type KafkaProducer struct {
	Topic    string
	Producer *commconf.KafkaProducer
}

type RPCClient struct {
	RegisterClient *commconf.RPCClient
	ManagerClient  *commconf.RPCClient
}

type MongoDB struct {
	*commconf.MongoDB
	OfflineMsgCollection string
}

// type ES struct {
// 	*commconf.ES
// 	Index string
// }

func init() {
	flag.StringVar(&confPath, "conf", "./logic.toml", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
