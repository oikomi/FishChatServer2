package dao

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/xmysql"
	"github.com/oikomi/FishChatServer2/server/register/conf"
	"golang.org/x/net/context"
)

const (
	_inUserSQL = "INSERT INTO user (uid, user_name, password) VALUES(?,?,?)"
)

type Mysql struct {
	DB         *xmysql.DB
	inUserStmt *xmysql.Stmt
}

func NewMysql() (mysql *Mysql) {
	mysql = &Mysql{
		DB: xmysql.NewMySQL(conf.Conf.Mysql.User),
	}
	mysql.initStmt()
	return
}

func (mysql *Mysql) initStmt() {
	mysql.inUserStmt = mysql.DB.Prepared(_inUserSQL)
}

func (mysql *Mysql) Insert(c context.Context, uid int64, userName, password string) (rows int64, err error) {
	res, err := mysql.inUserStmt.Exec(c, uid, userName, password)
	if err != nil {
		glog.Error(err)
		return
	}
	return res.RowsAffected()
}
