package main

import (
	// "context"
	// "encoding/json"
	"flag"
	// "github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/dao/kafka"
	"github.com/oikomi/FishChatServer2/job/msg/conf"
	// "time"
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
	// glog.Info(consumer)
	if consumer == nil {
		glog.Error("consumer is nil")
		return
	}

	for {
		msg, ok := <-consumer.ConsumerGroup.Messages()
		if !ok {
			glog.Info("consumeproc exit")
			return
		}
		if msg.Topic != conf.Conf.KafkaConsumer.Topics[0] {
			continue
		}
		glog.Info(string(msg.Value))
		// var mes string
		// if err = json.Unmarshal(msg.Value, mes); err != nil {
		// 	glog.Error("json.Unmarshal() error(%v)", err)
		// 	continue
		// }
	}

}
