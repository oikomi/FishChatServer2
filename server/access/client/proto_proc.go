package client

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/external"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/access/conf"
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
		UID:        reqLogin.UID,
		Token:      reqLogin.Token,
		AccessAddr: conf.Conf.RPCServer.Addr,
	}
	resLoginRPC, err := c.RPCClient.Logic.Login(reqLoginRPC)
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

func (c *Client) procReqPing(reqData []byte) (err error) {
	glog.Info("procReqPing")
	reqPing := &external.ReqPing{}
	if err = proto.Unmarshal(reqData, reqPing); err != nil {
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
	reqPingRPC := &rpc.PingReq{
		UID: reqPing.UID,
	}
	resPingRPC, err := c.RPCClient.Logic.Ping(reqPingRPC)
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
		Cmd:     external.PingCMD,
		ErrCode: resPingRPC.ErrCode,
		ErrStr:  resPingRPC.ErrStr,
	}); err != nil {
		glog.Error(err)
	}
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
		MsgType:   "p2p",
		MsgID:     reqSendP2PMsg.MsgID,
		Msg:       reqSendP2PMsg.Msg,
	}
	resSendP2PMsgRPC, err := c.RPCClient.Logic.SendP2PMsg(reqSendP2PMsgRPC)
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
	if err = c.Session.Send(&external.ResSendP2PMsg{
		Cmd:     external.SendP2PMsgCMD,
		ErrCode: resSendP2PMsgRPC.ErrCode,
		ErrStr:  resSendP2PMsgRPC.ErrStr,
	}); err != nil {
		glog.Error(err)
	}
	return
}

// func (c *Client) procAcceptP2PMsgAck(reqData []byte) (err error) {
// 	glog.Info("procAcceptP2PMsgAck")
// 	reqAcceptP2PMsgAck := &external.ReqAcceptP2PMsgAck{}
// 	if err = proto.Unmarshal(reqData, reqAcceptP2PMsgAck); err != nil {
// 		if err = c.Session.Send(&external.Error{
// 			Cmd:     external.ErrServerCMD,
// 			ErrCode: ecode.ServerErr.Uint32(),
// 			ErrStr:  ecode.ServerErr.String(),
// 		}); err != nil {
// 			glog.Error(err)
// 		}
// 		glog.Error(err)
// 		return
// 	}
// 	reqAcceptP2PMsgAckRPC := &rpc.AcceptP2PMsgAckReq{
// 		SourceUID: reqAcceptP2PMsgAck.SourceUID,
// 		TargetUID: reqAcceptP2PMsgAck.TargetUID,
// 		MsgID:     reqAcceptP2PMsgAck.MsgID,
// 	}
// 	// add rpc logic
// 	resAcceptP2PMsgAckRPC, err := c.RPCClient.Logic.AcceptP2PMsgAck(reqAcceptP2PMsgAckRPC)
// 	if err != nil {
// 		if err = c.Session.Send(&external.Error{
// 			Cmd:     external.ErrServerCMD,
// 			ErrCode: ecode.ServerErr.Uint32(),
// 			ErrStr:  ecode.ServerErr.String(),
// 		}); err != nil {
// 			glog.Error(err)
// 		}
// 		glog.Error(err)
// 		return
// 	}
// 	if err = c.Session.Send(&external.Error{
// 		Cmd:     external.SendP2PMsgCMD,
// 		ErrCode: resAcceptP2PMsgAckRPC.ErrCode,
// 		ErrStr:  resAcceptP2PMsgAckRPC.ErrStr,
// 	}); err != nil {
// 		glog.Error(err)
// 	}
// 	return
// }

func (c *Client) procSendGroupMsg(reqData []byte) (err error) {
	glog.Info("procSendGroupMsg")
	reqSendGroupMsg := &external.ReqSendGroupMsg{}
	if err = proto.Unmarshal(reqData, reqSendGroupMsg); err != nil {
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
	reqSendGroupMsgRPC := &rpc.SendGroupMsgReq{
		SourceUID: reqSendGroupMsg.SourceUID,
		MsgType:   "group",
		GroupID:   reqSendGroupMsg.GroupID,
		MsgID:     reqSendGroupMsg.MsgID,
		Msg:       reqSendGroupMsg.Msg,
	}
	// add rpc logic
	resSendGroupMsgRPC, err := c.RPCClient.Logic.SendGroupMsg(reqSendGroupMsgRPC)
	if err != nil {
		if err = c.Session.Send(&external.ResSendGroupMsg{
			Cmd:     external.ErrServerCMD,
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}); err != nil {
			glog.Error(err)
		}
		glog.Error(err)
		return
	}
	if err = c.Session.Send(&external.ResSendGroupMsg{
		Cmd:     external.SendGroupMsgCMD,
		ErrCode: resSendGroupMsgRPC.ErrCode,
		ErrStr:  resSendGroupMsgRPC.ErrStr,
	}); err != nil {
		glog.Error(err)
	}
	return
}

func (c *Client) procSyncMsg(reqData []byte) (err error) {
	glog.Info("procSyncMsg")
	reqSyncMsg := &external.ReqSyncMsg{}
	if err = proto.Unmarshal(reqData, reqSyncMsg); err != nil {
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
	reqSyncMsgRPC := &rpc.SyncMsgReq{
		UID:       reqSyncMsg.UID,
		CurrentID: reqSyncMsg.CurrentID,
	}
	resSyncMsgRPC, err := c.RPCClient.Logic.SyncMsg(reqSyncMsgRPC)
	if err != nil {
		if err = c.Session.Send(&external.ResSyncMsg{
			Cmd:     external.ErrServerCMD,
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}); err != nil {
			glog.Error(err)
		}
		glog.Error(err)
		return
	}
	tmpMsgs := make([]*external.OffsetMsg, 0)
	for _, v := range resSyncMsgRPC.Msgs {
		tmpMsg := &external.OffsetMsg{}
		tmpMsg.SourceUID = v.SourceUID
		tmpMsg.TargetUID = v.TargetUID
		tmpMsg.GroupID = v.GroupID
		tmpMsg.MsgType = v.MsgType
		tmpMsg.MsgID = v.MsgID
		tmpMsg.Msg = v.Msg
		tmpMsgs = append(tmpMsgs, tmpMsg)
	}
	if err = c.Session.Send(&external.ResSyncMsg{
		Cmd:       external.SyncMsgCMD,
		ErrCode:   resSyncMsgRPC.ErrCode,
		ErrStr:    resSyncMsgRPC.ErrStr,
		CurrentID: resSyncMsgRPC.CurrentID,
		Msgs:      tmpMsgs,
	}); err != nil {
		glog.Error(err)
	}
	return
}
