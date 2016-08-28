package main

import (
	// "encoding/binary"
	"flag"
	"github.com/golang/glog"
	// "fmt"
	// "encoding/binary"
	// "github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/libnet"
	// mybinary "github.com/oikomi/FishChatServer2/libnet/binary"
	// "fmt"
	// myproto "github.com/oikomi/FishChatServer2/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/codec"
	"github.com/oikomi/FishChatServer2/protocol"
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

func clientLoop(session *libnet.Session) {
	err := session.Send(&protocol.ReqMsgServer{
		Cmd: proto.Uint32(protocol.ReqMsgServerCMD),
	})
	// err = session.Send(&protocol.SelectMsgServerForClient{
	// 	Cmd:  proto.Uint32(1000001),
	// 	Addr: proto.String("122"),
	// })
	checkErr(err)
	// rsp, err := session.Receive()
	// checkErr(err)
	// glog.Info(rsp)
}

func main() {
	var addr string

	flag.StringVar(&addr, "addr", "127.0.0.1:17000", "echo server address")
	flag.Parse()
	protobuf := codec.Protobuf()
	// session, err := libnet.Connect("tcp", addr, libnet.Packet(2, 1024*1024, 1024, binary.BigEndian, TestCodec{}))
	client, err := libnet.Connect("tcp", addr, protobuf, 0)
	checkErr(err)
	clientLoop(client)
}
