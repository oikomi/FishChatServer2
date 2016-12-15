package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/codec"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/server/gateway/conf"
	"github.com/oikomi/FishChatServer2/server/gateway/job"
	"github.com/oikomi/FishChatServer2/server/gateway/server"
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
	gwServer := server.New()
	protobuf := codec.Protobuf()
	if gwServer.Server, err = libnet.Serve(conf.Conf.Server.Proto, conf.Conf.Server.Addr, protobuf, 0); err != nil {
		glog.Error(err)
		panic(err)
	}
	go job.ConfDiscoveryProc()
	gwServer.Loop()
}
