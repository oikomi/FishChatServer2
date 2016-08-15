package net

import (
	"net"
	"sync"
)

type Server struct {
	listener  net.Listener
	// About sessions
	maxSessionId uint64
	//sessionMaps  [sessionMapNum]sessionMap
	// About server start and stop
	stopOnce sync.Once
	stopWait sync.WaitGroup
	// Server state
	State interface{}
}

func NewServer(proto, addr string) (*Server, err error) {
	listener, err := net.Listen(proto, addr)
	server := &Server{
		listener:  listener,
	}
	return server, err
}