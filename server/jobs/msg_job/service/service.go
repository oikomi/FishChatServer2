package service

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/kafka"
	"github.com/oikomi/FishChatServer2/common/model"
	protoRPC "github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/jobs/msg_job/conf"
	"github.com/oikomi/FishChatServer2/server/jobs/msg_job/rpc"
	"sync"
	"time"
)

var (
	_module = "msg_job"
	_add    = "add"
)

type Service struct {
	c         *conf.Config
	waiter    *sync.WaitGroup
	consumer  *kafka.Consumer
	rpcClient *rpc.RPCClient
}

func New(c *conf.Config) (s *Service) {
	rpcClient, err := rpc.NewRPCClient()
	if err != nil {
		glog.Error(err)
		return
	}
	s = &Service{
		c:         c,
		waiter:    new(sync.WaitGroup),
		consumer:  kafka.NewConsumer(c.KafkaConsumer),
		rpcClient: rpcClient,
	}
	for s.consumer.ConsumerGroup == nil {
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
		msg, ok := <-s.consumer.ConsumerGroup.Messages()
		if !ok {
			glog.Error("consumeproc exit")
			return
		}
		glog.Info(string(msg.Value))
		if msg.Topic != s.c.KafkaConsumer.Topics[0] {
			continue
		}
		sendP2PMsgKafka := &model.SendP2PMsgKafka{}
		if err := json.Unmarshal(msg.Value, sendP2PMsgKafka); err != nil {
			glog.Error("json.Unmarshal() error ", err)
			continue
		}
		// if act.Action != _add {
		// 	continue
		// }
		// aid, err := strconv.ParseInt(string(msg.Key), 10, 64)
		// if err == nil {
		// 	ctx := context.Background()
		// 	now := time.Now()
		// 	if err := s.arcRPC.AddMoment(trace.NewContext(ctx, t), &model.ArgAid{Aid: aid}); err != nil {
		// 		if s.elk != nil {
		// 			tmsub := time.Now().Sub(now)
		// 			s.elk.Error(t.ID, "moment-job", err.Error(), tmsub.Nanoseconds())
		// 		}
		// 		log.Error("moment add(%d) error(%v)", aid, err)
		// 	}
		// }
		sendP2PMsgReq := &protoRPC.ASSendP2PMsgReq{
			SourceUID: sendP2PMsgKafka.UID,
			TargetUID: sendP2PMsgKafka.TargetUID,
			Msg:       sendP2PMsgKafka.Msg,
		}
		s.rpcClient.AccessServer.SendP2PMsg(sendP2PMsgReq)
		s.consumer.ConsumerGroup.CommitUpto(msg)
	}
}

func (s *Service) errproc() {
	errs := s.consumer.ConsumerGroup.Errors()
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
	return s.consumer.Close()
}

func (s *Service) Wait() {
	s.waiter.Wait()
}
