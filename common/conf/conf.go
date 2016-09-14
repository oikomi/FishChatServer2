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
