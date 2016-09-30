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
	producer := kafka.NewProducer(conf.Conf.KafkaProducer.Producer)
	glog.Info(producer)

	if err = producer.Input(context.Background(), &sarama.ProducerMessage{
		Topic: conf.Conf.KafkaProducer.Topic,
		Key:   sarama.StringEncoder("11"),
		Value: sarama.StringEncoder("22"),
	}); err != nil {
		glog.Error(err)
	}
	time.Sleep(time.Second)
}
