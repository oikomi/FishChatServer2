package server

import (
	"github.com/oikomi/FishChatServer2/gateway/client"
	"github.com/oikomi/FishChatServer2/gateway/conf"
	"github.com/oikomi/FishChatServer2/libnet"
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
		// glog.Info(req)

		// err = session.Send(&AddRsp{
		// 	req.(*AddReq).A + req.(*AddReq).B,
		// })
	}
}

func (s *Server) Loop() {
	for {
		session, err := s.Server.Accept()
		go s.sessionLoop(client.New(session))
	}
}
