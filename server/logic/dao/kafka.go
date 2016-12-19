package dao

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/kafka"
	"github.com/oikomi/FishChatServer2/common/model"
	"github.com/oikomi/FishChatServer2/server/logic/conf"
	"golang.org/x/net/context"
)

type KafkaProducer struct {
	producer         *kafka.Producer
	sendP2PMsgChan   chan *model.SendP2PMsgKafka
	sendGroupMsgChan chan *model.SendGroupMsgKafka
}

func NewKafkaProducer() (kafkaProducer *KafkaProducer) {
	producer := kafka.NewProducer(conf.Conf.KafkaProducer.Producer)
	kafkaProducer = &KafkaProducer{
		producer:         producer,
		sendP2PMsgChan:   make(chan *model.SendP2PMsgKafka, 1),
		sendGroupMsgChan: make(chan *model.SendGroupMsgKafka, 1),
	}
	return
}

func (kp *KafkaProducer) SendP2PMsg(data *model.SendP2PMsgKafka) {
	kp.sendP2PMsgChan <- data
}

func (kp *KafkaProducer) SendGroupMsg(data *model.SendGroupMsgKafka) {
	kp.sendGroupMsgChan <- data
}

func (kp *KafkaProducer) HandleSuccess() {
	var (
		pm *sarama.ProducerMessage
	)
	for {
		pm = <-kp.producer.Successes()
		if pm != nil {
			glog.Info("producer message success, partition:%d offset:%d key:%v valus:%s", pm.Partition, pm.Offset, pm.Key, pm.Value)
		}
	}
}

func (kp *KafkaProducer) HandleError() {
	var (
		err *sarama.ProducerError
	)
	for {
		err = <-kp.producer.Errors()
		if err != nil {
			glog.Error("producer message error, partition:%d offset:%d key:%v valus:%s error(%v)", err.Msg.Partition, err.Msg.Offset, err.Msg.Key, err.Msg.Value, err.Err)
		}
	}
}

func (kp *KafkaProducer) Process() {
	var sendP2PMsg *model.SendP2PMsgKafka
	var sendGroupMsg *model.SendGroupMsgKafka
	for {
		select {
		case sendP2PMsg = <-kp.sendP2PMsgChan:
			var (
				err    error
				vBytes []byte
			)
			if vBytes, err = json.Marshal(sendP2PMsg); err != nil {
				glog.Error(err)
				return
			}
			glog.Info("send to kafka : ", string(vBytes))
			if err := kp.producer.Input(context.Background(), &sarama.ProducerMessage{
				Topic: conf.Conf.KafkaProducer.P2PTopic,
				Key:   sarama.StringEncoder(model.SendP2PMsgKey),
				Value: sarama.ByteEncoder(vBytes),
			}); err != nil {
				glog.Error(err)
			}
		case sendGroupMsg = <-kp.sendGroupMsgChan:
			var (
				err    error
				vBytes []byte
			)
			if vBytes, err = json.Marshal(sendGroupMsg); err != nil {
				glog.Error(err)
				return
			}
			glog.Info("send to kafka : ", string(vBytes))
			if err := kp.producer.Input(context.Background(), &sarama.ProducerMessage{
				Topic: conf.Conf.KafkaProducer.GroupTopic,
				Key:   sarama.StringEncoder(model.SendGroupMsgKey),
				Value: sarama.ByteEncoder(vBytes),
			}); err != nil {
				glog.Error(err)
			}
		}
	}
}
