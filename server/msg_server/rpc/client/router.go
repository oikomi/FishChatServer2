package client

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/msg_server/conf"
	// "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RouterRPCCli struct {
	conn *grpc.ClientConn
}

func NewRouterRPCCli() (routerRPCCli *RouterRPCCli, err error) {
	conn, err := grpc.Dial(conf.Conf.RPCClient.ManagerClient.Addr, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
		return
	}
	routerRPCCli = &RouterRPCCli{
		conn: conn,
	}
	return
}

// func (routerRPCCli *RouterRPCCli) SendMsgP2P(targetUID, msg string) (err error) {
// 	r := rpc.NewRouterRPCClient(routerRPCCli.conn)
// 	res, err := r.SendMsgP2P(context.Background(), &rpc.SendMsgP2PReq{})
// 	if err != nil {
// 		glog.Error(err)
// 		return
// 	}
// 	glog.Info(res.ErrCode)
// 	return
// }
