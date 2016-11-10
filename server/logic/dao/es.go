package dao

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/server/logic/conf"
	elastic "gopkg.in/olivere/elastic.v5"
)

type ES struct {
	esCli *elastic.Client
}

func NewES() (es *ES, err error) {
	// client, err := elastic.NewClient(elastic.SetURL(conf.Conf.ES.ES.Addrs))
	client, err := elastic.NewClient(elastic.SetURL())
	if err != nil {
		glog.Error(err)
		return
	}
	es = &ES{
		esCli: client,
	}
	return
}
