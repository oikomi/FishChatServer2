package main

import (
	// "encoding/binary"
	"flag"
	"github.com/golang/glog"
	// "github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/codec"
	"github.com/oikomi/FishChatServer2/gateway/conf"
	"github.com/oikomi/FishChatServer2/libnet"
	// mybinary "github.com/oikomi/FishChatServer2/libnet/binary"
	// "github.com/oikomi/FishChatServer2/proto"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func checkErr(err error) {
	if err != nil {
		glog.Error(err)
	}
}

func sessionLoop(session *libnet.Session) {
	for {
		_, err := session.Receive()
		checkErr(err)
		// glog.Info(req)

		// err = session.Send(&AddRsp{
		// 	req.(*AddReq).A + req.(*AddReq).B,
		// })
	}
}

func serverLoop(server *libnet.Server) {
	for {
		session, err := server.Accept()
		checkErr(err)
		go sessionLoop(session)
	}
}

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		glog.Error("conf.Init() error: ", err)
		panic(err)
	}
	protobuf := codec.Protobuf()
	server, err := libnet.Serve(conf.Conf.Server.Proto, conf.Conf.Server.Addr, protobuf, 0 /* sync send */)
	if err != nil {
		glog.Error(err)
	}
	serverLoop(server)
}
