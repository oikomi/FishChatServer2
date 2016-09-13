package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/codec"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/server/manager/conf"
	"github.com/oikomi/FishChatServer2/server/manager/rpc"
	"github.com/oikomi/FishChatServer2/server/manager/server"
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
	manager := server.New(conf.Conf)
	protobuf := codec.Protobuf()
	manager.Server, err = libnet.Serve(conf.Conf.Server.Proto, conf.Conf.Server.Addr, protobuf, 0 /* sync send */)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	go rpc.RPCInit()
	rpc.NewMsgServerRPCCli(conf.Conf.RPCClient.MsgServerClient.Addr)
	manager.Loop()
}
