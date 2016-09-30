package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/conf"
	"github.com/wvanbergen/kafka/consumergroup"
	"golang.org/x/net/context"
	"time"
)

const (
	_module = "kafka"
)

var (
	ErrProducer = errors.New("kafka producer nil")
	ErrConsumer = errors.New("kafka consumer nil")
)

type Producer struct {
	sarama.AsyncProducer
	sarama.SyncProducer
	c   *conf.KafkaProducer
	env string
}

// NewProducer new kafka async producer and retry when has error.
func NewProducer(c *conf.KafkaProducer) (p *Producer) {
	var err error
	p = &Producer{
		c:   c,
		env: fmt.Sprintf("zookeeper%s@%v|brokers%v|sync(%t)", c.Zookeeper.Root, c.Zookeeper.Addrs, c.Brokers, c.Sync),
	}
	if !c.Sync {
		if err = p.asyncDial(); err != nil {
			go p.reAsyncDial()
		}
	} else {
		if err = p.syncDial(); err != nil {
			go p.reSyncDial()
		}
	}
	return
}

func (p *Producer) syncDial() (err error) {
	p.SyncProducer, err = sarama.NewSyncProducer(p.c.Brokers, nil)
	return
}

func (p *Producer) reSyncDial() {
	var err error
	for {
		if err = p.syncDial(); err == nil {
			glog.Info("kafka retry new sync producer ok")
			return
		} else {
			glog.Error("dial kafka producer error(%v)", err)
		}
		time.Sleep(time.Second)
	}
}

func (p *Producer) asyncDial() (err error) {
	if p.AsyncProducer, err = sarama.NewAsyncProducer(p.c.Brokers, nil); err == nil {
		go p.errproc()
		go p.successproc()
	}
	return
}

func (p *Producer) reAsyncDial() {
	var err error
	for {
		if err = p.asyncDial(); err == nil {
			glog.Info("kafka retry new async producer ok")
			return
		} else {
			glog.Error("dial kafka producer error(%v)", err)
		}
		time.Sleep(time.Second)
	}
	return
}

// errproc errors when aync producer publish messages.
// NOTE: Either Errors channel or Successes channel must be read. See the doc of AsyncProducer
func (p *Producer) errproc() {
	err := p.Errors()
	for {
		e, ok := <-err
		if !ok {
			return
		}
		glog.Error("kafka producer send message(%v) failed error(%v)", e.Msg, e.Err)
		// if c, ok := e.Msg.Metadata.(context.Context); ok {
		// 	// if t, ok := trace.FromContext(c); ok {
		// 	// 	t.ClientReceive()
		// 	// }
		// }
	}
}

func (p *Producer) successproc() {
	suc := p.Successes()
	for {
		msg, ok := <-suc
		if !ok {
			return
		}
		if _, ok := msg.Metadata.(context.Context); ok {
			// if t, ok := trace.FromContext(c); ok {
			// 	t.ClientReceive()
			// }
		}
	}
}

// Input send msg to kafka
// NOTE: If producer has beed created failed, the message will lose.
func (p *Producer) Input(c context.Context, msg *sarama.ProducerMessage) (err error) {
	if !p.c.Sync {
		if p.AsyncProducer == nil {
			err = ErrProducer
		} else {
			msg.Metadata = c
			// if t, ok := trace.FromContext(c); ok {
			// 	t = t.Fork()
			// 	t.ClientStart(_module, "async_input", p.env)
			// }
			p.AsyncProducer.Input() <- msg
		}
	} else {
		if p.SyncProducer == nil {
			err = ErrProducer
		} else {
			// if t, ok := trace.FromContext(c); ok {
			// 	t = t.Fork()
			// 	t.ClientStart(_module, "sync_input", p.env)
			// 	defer t.ClientReceive()
			// }
			_, _, err = p.SyncProducer.SendMessage(msg)
		}
	}
	return
}

func (p *Producer) Close() (err error) {
	if !p.c.Sync {
		if p.AsyncProducer != nil {
			return p.AsyncProducer.Close()
		}
	}
	if p.SyncProducer != nil {
		return p.SyncProducer.Close()
	}
	return
}

// kafka consumer
type Consumer struct {
	ConsumerGroup *consumergroup.ConsumerGroup
	c             *conf.KafkaConsumer
}

func NewConsumer(c *conf.KafkaConsumer) (kc *Consumer) {
	var err error
	kc = &Consumer{
		c: c,
	}
	if err = kc.dial(); err != nil {
		go kc.redial()
	}
	return
}

func (c *Consumer) dial() (err error) {
	cfg := consumergroup.NewConfig()
	if c.c.Offset {
		cfg.Offsets.Initial = sarama.OffsetNewest
	} else {
		cfg.Offsets.Initial = sarama.OffsetOldest
	}
	cfg.Zookeeper.Chroot = c.c.Zookeeper.Root
	cfg.Zookeeper.Timeout = time.Duration(c.c.Zookeeper.Timeout)
	c.ConsumerGroup, err = consumergroup.JoinConsumerGroup(c.c.Group, c.c.Topics, c.c.Zookeeper.Addrs, cfg)
	return
}

func (c *Consumer) redial() {
	var err error
	for {
		if err = c.dial(); err == nil {
			glog.Info("kafka retry new consumer ok")
			return
		} else {
			glog.Error("dial kafka consumer error ", err)
		}
		time.Sleep(time.Second)
	}
}

func (c *Consumer) Close() error {
	if c.ConsumerGroup != nil {
		return c.ConsumerGroup.Close()
	}
	return nil
}
