package server

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/codec"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol"
	"github.com/oikomi/FishChatServer2/server/gateway/client"
	"github.com/oikomi/FishChatServer2/server/gateway/conf"
	// "github.com/oikomi/FishChatServer2/server/gateway/proto_proc"
)

type Server struct {
	Config    *conf.Config
	Server    *libnet.Server
	ProtoProc *ProtoProc
}

func New(config *conf.Config) (s *Server) {
	s = &Server{
		Config:    config,
		ProtoProc: NewProtoProc(config),
	}
	return
}

func (s *Server) sessionLoop(client *client.Client) {
	// protoProc := proto_proc.New(s)
	for {
		reqData, err := client.Session.Receive()
		if err != nil {
			glog.Error(err)
		}
		if reqData != nil {
			baseCMD := &protocol.Base{}
			if err = proto.Unmarshal(reqData, baseCMD); err != nil {
				err = client.Session.Send(&protocol.Error{
					ErrCode: proto.Uint32(ecode.ServerErr.Uint32()),
					ErrStr:  proto.String(ecode.ServerErr.String()),
				})
			}
			s.ProtoProc.Parse(baseCMD.GetCmd(), reqData, client)
		}
	}
}

func (s *Server) Loop() {
	for {
		session, err := s.Server.Accept()
		if err != nil {
			glog.Error(err)
		}
		go s.sessionLoop(client.New(session))
	}
}
