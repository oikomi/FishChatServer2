package client

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/external"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/access/global"
)

func (c *Client) procReqLogin(reqData []byte) (err error) {
	glog.Info("procReqLogin")
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
		Cmd:     external.LoginCMD,
		ErrCode: resLoginRPC.ErrCode,
		ErrStr:  resLoginRPC.ErrStr,
	}); err != nil {
		glog.Error(err)
	}
	// success
	global.GSessions[reqLogin.UID] = c.Session
	return
}

func (c *Client) procSendP2PMsg(reqData []byte) (err error) {
	glog.Info("procSendP2PMsg")
	reqSendP2PMsg := &external.ReqSendP2PMsg{}
	if err = proto.Unmarshal(reqData, reqSendP2PMsg); err != nil {
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
	reqSendP2PMsgRPC := &rpc.SendP2PMsgReq{
		SourceUID: reqSendP2PMsg.SourceUID,
		TargetUID: reqSendP2PMsg.TargetUID,
		Msg:       reqSendP2PMsg.Msg,
	}
	resSendP2PMsgRPC, err := c.RPCClient.MsgServer.SendP2PMsg(reqSendP2PMsgRPC)
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
		Cmd:     external.SendP2PMsgCMD,
		ErrCode: resSendP2PMsgRPC.ErrCode,
		ErrStr:  resSendP2PMsgRPC.ErrStr,
	}); err != nil {
		glog.Error(err)
	}
	return
}
