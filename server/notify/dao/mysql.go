package dao

import (
	"database/sql"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/xmysql"
	"github.com/oikomi/FishChatServer2/server/notify/conf"
	"github.com/oikomi/FishChatServer2/server/notify/model"
	"golang.org/x/net/context"
)

const (
	_upUserMsgIDSQL  = "UPDATE user_msg_id SET current_msg_id=?,total_msg_id=? WHERE uid=?"
	_getUserMsgIDSQL = "SELECT uid,current_msg_id,total_msg_id FROM user_msg_id WHERE uid=?"
)

type Mysql struct {
	im               *xmysql.DB
	upUserMsgIDStmt  *xmysql.Stmt
	getUserMsgIDStmt *xmysql.Stmt
}

func NewMysql() (mysql *Mysql) {
	mysql = &Mysql{
		im: xmysql.NewMySQL(conf.Conf.Mysql.IM),
	}
	mysql.initStmt()
	return
}

func (mysql *Mysql) initStmt() {
	mysql.upUserMsgIDStmt = mysql.im.Prepared(_upUserMsgIDSQL)
	mysql.getUserMsgIDStmt = mysql.im.Prepared(_getUserMsgIDSQL)
}

func (mysql *Mysql) UpdateUserMsgID(c context.Context, uid, currentMsgID, totalMsgID int64) (rows int64, err error) {
	res, err := mysql.upUserMsgIDStmt.Exec(c, currentMsgID, totalMsgID, uid)
	if err != nil {
		glog.Error(err)
		return
	}
	return res.RowsAffected()
}

func (mysql *Mysql) GetUserMsgID(c context.Context, uid int64) (userMsgID *model.UserMsgID, err error) {
	row := mysql.im.QueryRow(c, _getUserMsgIDSQL, uid)
	userMsgID = &model.UserMsgID{}
	if err = row.Scan(&userMsgID.UID, &userMsgID.CurrentMsgID, &userMsgID.TotalMsgID); err != nil {
		if err == sql.ErrNoRows {
			userMsgID = nil
			err = nil
		} else {
			glog.Error(err)
		}
	}
	return
}
