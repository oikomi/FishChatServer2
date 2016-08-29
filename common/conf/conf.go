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

type Zookeeper struct {
	Root    string
	Addrs   []string
	Timeout time.Duration
}
