package main

import (
	"io"
	"flag"
	"github.com/oikomi/FishChatServer2/gateway/conf"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/net"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		glog.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	listener, err := net.NewServer(conf.Conf.Server.Proto, conf.Conf.Server.Addr)
	if err != nil {
		panic(err)
	}
	println("server start:", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		go io.Copy(conn, conn)
	}
}
