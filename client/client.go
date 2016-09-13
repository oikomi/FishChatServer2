package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/codec"
	"github.com/oikomi/FishChatServer2/libnet"
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
		Cmd: protocol.ReqMsgServerCMD,
	})
	checkErr(err)
	rsp, err := session.Receive()
	checkErr(err)
	glog.Info(string(rsp))
	if rsp != nil {
		baseCMD := &protocol.Base{}
		if err = proto.Unmarshal(rsp, baseCMD); err != nil {

		}
		switch baseCMD.Cmd {
		case protocol.ResSelectMsgServerForClientCMD:
			resSelectMsgServerForClientPB := &protocol.ResSelectMsgServerForClient{}
			proto.Unmarshal(rsp, resSelectMsgServerForClientPB)
			glog.Info(resSelectMsgServerForClientPB.Addr)
		}

	}
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:11000", "echo server address")
	flag.Parse()
	protobuf := codec.Protobuf()
	// session, err := libnet.Connect("tcp", addr, libnet.Packet(2, 1024*1024, 1024, binary.BigEndian, TestCodec{}))
	client, err := libnet.Connect("tcp", addr, protobuf, 0)
	checkErr(err)
	clientLoop(client)
}
