package service

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/kafka"
	"github.com/oikomi/FishChatServer2/jobs/msg_job/conf"
	"sync"
	"time"
)

var (
	_module = "msg_job"
	_add    = "add"
)

// action the message struct of kafka
type action struct {
	Action string `json:"action"`
}

type Service struct {
	c        *conf.Config
	waiter   *sync.WaitGroup
	consumer *kafka.Consumer
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		waiter:   new(sync.WaitGroup),
		consumer: kafka.NewConsumer(c.KafkaConsumer),
	}
	for s.consumer.ConsumerGroup == nil {
		time.Sleep(time.Second)
	}
	for i := 0; i < 8; i++ {
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
		msg, ok := <-s.consumer.ConsumerGroup.Messages()
		if !ok {
			glog.Info("consumeproc exit")
			return
		}
		if msg.Topic != s.c.KafkaConsumer.Topics[0] {
			continue
		}
		act := &action{}
		if err := json.Unmarshal(msg.Value, act); err != nil {
			glog.Error("json.Unmarshal() error ", err)
			continue
		}
		if act.Action != _add {
			continue
		}
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
		// glog.Error("topic(%s) partition(%d) error(%v)", err.Topic, err.Partition, err.Err)
	}
}

func (s *Service) Close() error {
	return s.consumer.Close()
}

func (s *Service) Wait() {
	s.waiter.Wait()
}
