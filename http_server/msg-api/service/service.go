package service

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/http_server/msg-api/model"
	managerRpc "github.com/oikomi/FishChatServer2/http_server/msg-api/rpc"
	// "github.com/oikomi/FishChatServer2/protocol/rpc"
)

type Service struct {
	rpcClient *managerRpc.RPCClient
}

func New() (service *Service, err error) {
	rpcClient, err := managerRpc.NewRPCClient()
	if err != nil {
		glog.Error(err)
		return
	}
	service = &Service{
		rpcClient: rpcClient,
	}
	return
}

// func (s *Service) GetOfflineMsgs(uid int64) (offlineMsgsModel *model.OfflineMsgs, err error) {
// 	rgGetOfflineMsg := &rpc.MGOfflineMsgReq{
// 		Uid: uid,
// 	}
// 	res, err := s.rpcClient.Manager.GetOfflineMsgs(rgGetOfflineMsg)
// 	if err != nil {
// 		glog.Error(err)
// 		return
// 	}
// 	offlineMsgsModel = &model.OfflineMsgs{}
// 	for _, msg := range res.Msgs {
// 		tmpMsg := &model.OfflineMsg{
// 			SourceUID: msg.SourceUID,
// 			TargetUID: msg.TargetUID,
// 			MsgID:     msg.MsgID,
// 			Msg:       msg.Msg,
// 		}
// 		offlineMsgsModel.Msgs = append(offlineMsgsModel.Msgs, tmpMsg)
// 	}
// 	return
// }
