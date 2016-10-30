package conf

import (
	"flag"
	// "go-common/business/service/identify"
	"github.com/oikomi/FishChatServer2/common/conf"
	// "go-common/xlog"
	"github.com/BurntSushi/toml"
	"github.com/oikomi/FishChatServer2/common/xtime"
)

var (
	confPath string
	Conf     *Config
)

type Config struct {
	conf.CommConf
	// http
	MultiHTTP *conf.MultiHTTP
	// redis
	Redis *Redis
	// kafka
	Kafka *Kafka
	// ELK
	// ELK *xlog.ELKConfig
	// // tracer
	// Tracer *conf.Tracer
	// // hbase
	// HBase *conf.HBase
}

type Redis struct {
	*conf.Redis
	Expire xtime.Duration
}

// type RPC struct {
// 	Account *conf.RPCClient2
// 	Archive *conf.RPCClient2
// }

type Kafka struct {
	Topic      string
	FirstTopic string
	Producer   *conf.KafkaProducer
}

func init() {
	flag.StringVar(&confPath, "conf", "./history-example.toml", "config path")
}

func Init() error {
	if _, err := toml.DecodeFile(confPath, &Conf); err != nil {
		return err
	}
	return nil
}
