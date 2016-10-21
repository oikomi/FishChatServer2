package main

import (
	"flag"
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/codec"
	// "github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/server/msg_server/conf"
	"github.com/oikomi/FishChatServer2/server/msg_server/rpc"
	// "github.com/oikomi/FishChatServer2/server/msg_server/server"
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
	rpc.RPCServerInit()
}
