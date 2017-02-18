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
	configFile             string
	Server                 *commconf.Server
	ServiceDiscoveryServer *commconf.ServiceDiscoveryServer
	RPCServer              *commconf.RPCServer
	RPCClient              *RPCClient
	Etcd                   *commconf.Etcd
	KafkaProducer          *KafkaProducer
}

type KafkaProducer struct {
	P2PTopic   string
	GroupTopic string
	Producer   *commconf.KafkaProducer
}

type RPCClient struct {
	RegisterClient *commconf.ServiceDiscoveryClient
	ManagerClient  *commconf.ServiceDiscoveryClient
	IdgenClient    *commconf.ServiceDiscoveryClient
	NotifyClient   *commconf.ServiceDiscoveryClient
}

// type MongoDB struct {
// 	*commconf.MongoDB
// 	OfflineMsgCollection string
// }

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
