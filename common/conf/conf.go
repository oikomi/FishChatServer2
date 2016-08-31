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
