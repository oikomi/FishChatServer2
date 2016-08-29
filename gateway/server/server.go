package server

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/codec"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/gateway/client"
	"github.com/oikomi/FishChatServer2/gateway/conf"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol"
)

type Server struct {
	Config *conf.Config
	Server *libnet.Server
}

func New(config *conf.Config) (s *Server) {
	s = &Server{
		Config: config,
	}
	return
}

func (s *Server) sessionLoop(client *client.Client) {
	for {
		reqData, err := client.Session.Receive()
		if err != nil {
			glog.Error(err)
		}
		if reqData != nil {
			// cmd := codec.GetUint32BE(reqData[:4])
			// glog.Info(cmd)
			test := &protocol.Base{}
			err := proto.Unmarshal(reqData, test)
			if err != nil {

			}
			glog.Info(test.GetCmd())
		}
		// glog.Info(req)

		// err = session.Send(&AddRsp{
		// 	req.(*AddReq).A + req.(*AddReq).B,
		// })
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
