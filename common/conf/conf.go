package conf

import (
	"github.com/oikomi/FishChatServer2/common/xtime"
)

type CommConf struct {
	Ver     string
	LogPath string
}

// =================================== HTTP ==================================
// HTTPServer http server settings.
type HTTPServer struct {
	Addrs        []string
	MaxListen    int32
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

// HTTPClient http client settings.
type HTTPClient struct {
	Dial      xtime.Duration
	Timeout   xtime.Duration
	KeepAlive xtime.Duration
	Timer     int
}

// MultiHttp outer/inner/local http server settings.
type MultiHTTP struct {
	Outer *HTTPServer
	Inner *HTTPServer
	Local *HTTPServer
}

type Server struct {
	Proto string
	Addr  string
}

type RPCServer struct {
	Proto string
	Addr  string
}

type ConfDiscovery struct {
	Role     string
	Interval xtime.Duration
}

type ServiceDiscoveryServer struct {
	ServiceName string
	RPCAddr     string
	EtcdAddr    string
	Interval    xtime.Duration
	TTL         xtime.Duration
}

type ServiceDiscoveryClient struct {
	ServiceName string
	EtcdAddr    string
	Balancer    string
}

type Etcd struct {
	Name    string
	Root    string
	Addrs   []string
	Timeout xtime.Duration
}

type Zookeeper struct {
	Root    string
	Addrs   []string
	Timeout xtime.Duration
}

// Redis client settings.
type Redis struct {
	Name         string // redis name, for trace
	Proto        string
	Addr         string
	Active       int // pool
	Idle         int // pool
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
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

type MySQL struct {
	Name   string // for trace
	DSN    string // data source name
	Active int    // pool
	Idle   int    // pool
}

type MongoDB struct {
	Addrs       string
	DB          string
	DialTimeout xtime.Duration
}

type ES struct {
	Addrs string
}
