package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/codec"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/server/msg_server/conf"
	"github.com/oikomi/FishChatServer2/server/msg_server/rpc"
	"github.com/oikomi/FishChatServer2/server/msg_server/server"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func main() {
	var err error
	flag.Parse()
	if err = conf.Init(); err != nil {
		glog.Error("conf.Init() error: ", err)
		panic(err)
	}
	msgServer := server.New()
	protobuf := codec.Protobuf()
	msgServer.Server, err = libnet.Serve(conf.Conf.Server.Proto, conf.Conf.Server.Addr, protobuf, 0)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	rpcClient, err := rpc.NewRPCClient()
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	msgServer.SDHeart()
	go rpc.RPCServerInit()
	msgServer.Loop(rpcClient)
}
