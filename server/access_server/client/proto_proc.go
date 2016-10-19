package client

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/external"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
)

func (c *Client) procReqLogin(reqData []byte) (err error) {
	reqLogin := &external.ReqLogin{}
	if err = proto.Unmarshal(reqData, reqLogin); err != nil {
		if err = c.Session.Send(&external.Error{
			Cmd:     external.ErrServerCMD,
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}); err != nil {
			glog.Error(err)
		}
		return
	}
	reqLoginRPC := &rpc.LoginReq{
		UID:   reqLogin.UID,
		Token: reqLogin.Token,
	}
	resLoginRPC, err := c.RPCClient.MsgServer.Login(reqLoginRPC)
	if err != nil {
		if err = c.Session.Send(&external.Error{
			Cmd:     external.ErrServerCMD,
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}); err != nil {
			glog.Error(err)
		}
		glog.Error(err)
		return
	}
	if err = c.Session.Send(&external.Error{
		Cmd:     external.ReqLoginCMD,
		ErrCode: resLoginRPC.ErrCode,
		ErrStr:  resLoginRPC.ErrStr,
	}); err != nil {
		glog.Error(err)
	}

	// // FIXME
	// if reqLogin.Token == "" || reqLogin.UID < 0 {
	// 	if err = c.Session.Send(&external.ResLogin{
	// 		Cmd:     external.ReqLoginCMD,
	// 		ErrCode: ecode.NoToken.Uint32(),
	// 		ErrStr:  ecode.NoToken.String(),
	// 	}); err != nil {
	// 		glog.Error(err)
	// 	}
	// 	return
	// }
	// // success
	// if err = c.Session.Send(&external.ResLogin{
	// 	Cmd:     external.ResSelectAccessServerForClientCMD,
	// 	ErrCode: ecode.OK.Uint32(),
	// 	ErrStr:  ecode.OK.String(),
	// }); err != nil {
	// 	glog.Error(err)
	// }
	return
}

func (c *Client) procSendP2PMsg(reqData []byte) (err error) {
	reqSendP2PMsg := &external.ReqSendP2PMsg{}
	if err = proto.Unmarshal(reqData, reqSendP2PMsg); err != nil {
		if err = c.Session.Send(&external.Error{
			Cmd:     external.ErrServerCMD,
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}); err != nil {
			glog.Error(err)
		}
		return
	}
	return
}
