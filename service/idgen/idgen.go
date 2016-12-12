package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/service/idgen/conf"
	"github.com/oikomi/FishChatServer2/service/idgen/rpc"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		glog.Error("conf.Init() error: ", err)
		panic(err)
	}
	rpc.RPCServerInit()
}
