package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/codec"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/server/router/conf"
	"github.com/oikomi/FishChatServer2/server/router/rpc"
	"github.com/oikomi/FishChatServer2/server/router/server"
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
	router := server.New(conf.Conf)
	protobuf := codec.Protobuf()
	router.Server, err = libnet.Serve(conf.Conf.Server.Proto, conf.Conf.Server.Addr, protobuf, 0 /* sync send */)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	// monitor.SDHeart()
	go rpc.RPCInit()
	rpc.NewMsgServerRPCCli(conf.Conf.RPCClient.MsgServerClient.Addr)
	router.Loop()
}
