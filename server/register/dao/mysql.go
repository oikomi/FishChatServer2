package dao

import (
	"github.com/oikomi/FishChatServer2/common/dao/xmysql"
	"github.com/oikomi/FishChatServer2/server/register/conf"
)

type Mysql struct {
	DB *xmysql.DB
}

func NewMysql() (mysql *Mysql) {
	mysql = &Mysql{
		DB: xmysql.NewMySQL(conf.Conf.Mysql),
	}
	return
}
