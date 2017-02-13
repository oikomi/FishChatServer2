package dao

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/xhbase"
	"github.com/oikomi/FishChatServer2/server/manager/conf"
	"github.com/tsuna/gohbase/hrpc"
	"golang.org/x/net/context"
)

type HBase struct {
	client *xhbase.Client
}

func NewHBase() *HBase {
	return &HBase{
		client: xhbase.NewClient(conf.Conf.HBase.ZKAddr),
	}
}

func (h *HBase) GetMsgs(ctx context.Context, rowKey string) (res *hrpc.Result, err error) {
	getRequest, err := hrpc.NewGetStr(ctx, conf.Conf.HBase.Table, rowKey)
	if err != nil {
		glog.Error(err)
		return
	}
	res, err = h.client.Get(ctx, getRequest)
	if err != nil {
		glog.Error(err)
	}
	return
}
