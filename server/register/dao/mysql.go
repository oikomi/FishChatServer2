package dao

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/xmysql"
	"github.com/oikomi/FishChatServer2/server/register/conf"
	"golang.org/x/net/context"
)

const (
	_inUserSQL  = "INSERT INTO user (uid, user_name, password) VALUES(?,?,?)"
	_inGroupSQL = "INSERT INTO group (gid, group_name, owner_id) VALUES(?,?,?)"
)

type Mysql struct {
	im          *xmysql.DB
	inUserStmt  *xmysql.Stmt
	inGroupStmt *xmysql.Stmt
}

func NewMysql() (mysql *Mysql) {
	mysql = &Mysql{
		im: xmysql.NewMySQL(conf.Conf.Mysql.IM),
	}
	mysql.initStmt()
	return
}

func (mysql *Mysql) initStmt() {
	mysql.inUserStmt = mysql.im.Prepared(_inUserSQL)
	mysql.inGroupStmt = mysql.im.Prepared(_inGroupSQL)
}

func (mysql *Mysql) InsertUser(c context.Context, uid int64, userName, password string) (rows int64, err error) {
	res, err := mysql.inUserStmt.Exec(c, uid, userName, password)
	if err != nil {
		glog.Error(err)
		return
	}
	return res.RowsAffected()
}

func (mysql *Mysql) InsertGroup(c context.Context, gid, ownerID int64, groupName string) (rows int64, err error) {
	res, err := mysql.inGroupStmt.Exec(c, gid, groupName, ownerID)
	if err != nil {
		glog.Error(err)
		return
	}
	return res.RowsAffected()
}
