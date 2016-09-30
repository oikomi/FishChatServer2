package main

import (
	"context"
	"flag"
	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/dao/kafka"
	"github.com/oikomi/FishChatServer2/job/msg/conf"
	"time"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func main() {
	var err error
	flag.Parse()
	if err = conf.Init(); err != nil {
		glog.Error("conf.Init() error: ", err)
		panic(err)
	}
	consumer := kafka.NewConsumer(conf.Conf.KafkaConsumer)
	glog.Info(consumer)

}
