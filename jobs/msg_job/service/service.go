package service

import (
	"encoding/json"
	"github.com/golang/glog"
	commmodel "github.com/oikomi/FishChatServer2/common/model"
	"github.com/oikomi/FishChatServer2/jobs/msg_job/conf"
	"github.com/oikomi/FishChatServer2/jobs/msg_job/dao"
	"github.com/oikomi/FishChatServer2/jobs/msg_job/rpc"
	protoRPC "github.com/oikomi/FishChatServer2/protocol/rpc"
	"sync"
	"time"
)

var (
	_module = "msg_job"
)

type Service struct {
	c         *conf.Config
	waiter    *sync.WaitGroup
	dao       *dao.Dao
	rpcClient *rpc.RPCClient
}

func New(c *conf.Config) (s *Service) {
	rpcClient, err := rpc.NewRPCClient()
	if err != nil {
		glog.Error(err)
		return
	}
	dao, err := dao.NewDao()
	if err != nil {
		glog.Error(err)
		return
	}
	s = &Service{
		c:         c,
		waiter:    new(sync.WaitGroup),
		dao:       dao,
		rpcClient: rpcClient,
	}
	for s.dao.Kafka.Consumer.ConsumerGroup == nil {
		time.Sleep(time.Second)
	}
	for i := 0; i < 1; i++ {
		glog.Info("start proc")
		go s.consumeproc()
	}
	go s.errproc()
	return
}

func (s *Service) consumeproc() {
	s.waiter.Add(1)
	defer s.waiter.Done()
	for {
		glog.Info("start consume...")
		msg, ok := <-s.dao.Kafka.Consumer.ConsumerGroup.Messages()
		if !ok {
			glog.Error("consumeproc exit")
			return
		}
		glog.Info(string(msg.Value))
		if msg.Topic != s.c.KafkaConsumer.Topics[0] {
			continue
		}
		sendP2PMsgKafka := &commmodel.SendP2PMsgKafka{}
		if err := json.Unmarshal(msg.Value, sendP2PMsgKafka); err != nil {
			glog.Error("json.Unmarshal() error ", err)
			continue
		}
		if sendP2PMsgKafka.Online {
			sendP2PMsgReq := &protoRPC.ASSendP2PMsgFromJobReq{
				SourceUID:        sendP2PMsgKafka.SourceUID,
				TargetUID:        sendP2PMsgKafka.TargetUID,
				MsgID:            sendP2PMsgKafka.MsgID,
				Msg:              sendP2PMsgKafka.Msg,
				AccessServerAddr: sendP2PMsgKafka.AccessServerAddr,
			}
			_, err := s.rpcClient.AccessServer.SendP2PMsgFromJob(sendP2PMsgReq)
			if err != nil {
				// store offline msg
				glog.Error(err)
			}
		} else {
			// set offline msg
			offlineMsg := &commmodel.OfflineMsg{
				MsgID:     sendP2PMsgKafka.MsgID,
				SourceUID: sendP2PMsgKafka.SourceUID,
				TargetUID: sendP2PMsgKafka.TargetUID,
				Msg:       sendP2PMsgKafka.Msg,
			}
			if err := s.dao.MongoDB.StoreOfflineMsg(offlineMsg); err != nil {
				glog.Error(err)
				return
			}
		}
		s.dao.Kafka.Consumer.ConsumerGroup.CommitUpto(msg)
	}
}

func (s *Service) errproc() {
	errs := s.dao.Kafka.Consumer.ConsumerGroup.Errors()
	for {
		err, ok := <-errs
		if !ok {
			glog.Info("errproc exit")
			return
		}
		glog.Error(err)
	}
}

func (s *Service) Close() error {
	return s.dao.Kafka.Consumer.Close()
}

func (s *Service) Wait() {
	s.waiter.Wait()
}
