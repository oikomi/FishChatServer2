package main

import (
	// "encoding/binary"
	"flag"
	"github.com/golang/glog"
	// "github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/gateway/conf"
	// "github.com/oikomi/FishChatServer2/libnet"
	// mybinary "github.com/oikomi/FishChatServer2/libnet/binary"
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

}
