package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/codec"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/server/msg_server/conf"
	"github.com/oikomi/FishChatServer2/server/msg_server/pb"
	"github.com/oikomi/FishChatServer2/server/msg_server/server"
	"google.golang.org/grpc"
	"net"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func rpcInit() {
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error("failed to listen: %v", err)
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	s.Serve(lis)
}

func main() {
	var err error
	flag.Parse()
	if err = conf.Init(); err != nil {
		glog.Error("conf.Init() error: ", err)
		panic(err)
	}
	msgServer := server.New(conf.Conf)
	protobuf := codec.Protobuf()
	msgServer.Server, err = libnet.Serve(conf.Conf.Server.Proto, conf.Conf.Server.Addr, protobuf, 0 /* sync send */)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	msgServer.SDHeart()
	msgServer.Loop()
	rpcInit()
}
