package dao

import (
	"database/sql"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/xmysql"
	"github.com/oikomi/FishChatServer2/server/register/conf"
	"github.com/oikomi/FishChatServer2/server/register/model"
	"golang.org/x/net/context"
)

const (
	_getUserSQL           = "SELECT uid, user_name, password FROM user WHERE uid=?"
	_inUserSQL            = "INSERT INTO user (uid, user_name, password) VALUES(?,?,?)"
	_inGroupSQL           = "INSERT INTO `group` (gid, group_name, owner_id) VALUES(?,?,?)"
	_inUserMsgIDSQL       = "INSERT INTO user_msg_id (uid, current_msg_id, total_msg_id) VALUES(?,?,?)"
	_inUserGroupSQL       = "INSERT INTO `user_group` (gid, uid) VALUES(?,?)"
	_getUsersByGroupIDSQL = "SELECT uid FROM user_group WHERE gid=?"
)

type Mysql struct {
	im                *xmysql.DB
	getUserStmt       *xmysql.Stmt
	inUserStmt        *xmysql.Stmt
	inGroupStmt       *xmysql.Stmt
	inUserMsgIDStmt   *xmysql.Stmt
	inUserGroupStmt   *xmysql.Stmt
	getUsersByGroupID *xmysql.Stmt
}

func NewMysql() (mysql *Mysql) {
	mysql = &Mysql{
		im: xmysql.NewMySQL(conf.Conf.Mysql.IM),
	}
	mysql.initStmt()
	return
}

func (mysql *Mysql) initStmt() {
	mysql.getUserStmt = mysql.im.Prepared(_getUserSQL)
	mysql.inUserStmt = mysql.im.Prepared(_inUserSQL)
	mysql.inGroupStmt = mysql.im.Prepared(_inGroupSQL)
	mysql.inUserMsgIDStmt = mysql.im.Prepared(_inUserMsgIDSQL)
	mysql.inUserGroupStmt = mysql.im.Prepared(_inUserGroupSQL)
	mysql.getUsersByGroupID = mysql.im.Prepared(_getUsersByGroupIDSQL)
}

func (mysql *Mysql) GetUser(c context.Context, uid int64) (res *model.User, err error) {
	res = &model.User{}
	row, err := mysql.getUserStmt.QueryRow(c, uid)
	if err != nil {
		glog.Error(err)
		return
	}
	if err = row.Scan(&res.Uid, &res.UserName, &res.Password); err != nil {
		if err == sql.ErrNoRows {
			res = nil
			// err = nil
		} else {
			glog.Error(err)
		}
	}
	return
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

func (mysql *Mysql) InsertUserMsgID(c context.Context, uid, currentMsgID, totalMsgID int64) (rows int64, err error) {
	res, err := mysql.inUserMsgIDStmt.Exec(c, uid, currentMsgID, totalMsgID)
	if err != nil {
		glog.Error(err)
		return
	}
	return res.RowsAffected()
}

func (mysql *Mysql) InsertUserGroup(c context.Context, uid, gid int64) (rows int64, err error) {
	res, err := mysql.inUserGroupStmt.Exec(c, gid, uid)
	if err != nil {
		glog.Error(err)
		return
	}
	return res.RowsAffected()
}

func (mysql *Mysql) GetUsersByGroupID(c context.Context, gid int64) (uids []int64, err error) {
	rows, err := mysql.getUsersByGroupID.Query(c, gid)
	if err != nil {
		glog.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uid int64
		if err = rows.Scan(&uid); err != nil {
			glog.Error(err)
			return
		}
		uids = append(uids, uid)
	}
	return
}
