package conf

import (
	"time"
)

type CommConf struct {
	Ver     string
	LogPath string
}

type Server struct {
	Proto string
	Addr  string
}

type RPCServer struct {
	Proto string
	Addr  string
}

type RPCClient struct {
	Addr string
}

type Etcd struct {
	Name    string
	Root    string
	Addrs   []string
	Timeout time.Duration
}

type Zookeeper struct {
	Root    string
	Addrs   []string
	Timeout time.Duration
}

type Redis struct {
	Name         string
	Proto        string
	Addr         string
	Active       int
	Idle         int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// KafkaProducer kafka producer settings.
type KafkaProducer struct {
	Zookeeper *Zookeeper
	Brokers   []string
	Sync      bool // true: sync, false: async
}

// KafkaConsumer kafka client settings.
type KafkaConsumer struct {
	Group     string
	Topics    []string
	Offset    bool // true: new, false: old
	Zookeeper *Zookeeper
}
