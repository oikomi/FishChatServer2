package conf

import (
	"github.com/oikomi/FishChatServer2/common/xtime"
	// "time"
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

type RPCClient struct {
	Addr string
}

type ServiceDiscovery struct {
	Role     string
	Interval xtime.Duration
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

type MongoDB struct {
	Addrs       string
	DB          string
	DialTimeout xtime.Duration
}
