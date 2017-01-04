package xmysql

import (
	"database/sql"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/conf"
)

type DB struct {
	*sql.DB
	env string
}

func NewMySQL(c *conf.MySQL) (db *DB) {
	var err error
	if db, err = Open("mysql", c.DSN); err != nil {
		glog.Error(err)
		panic(err)
	}
	db.env = c.Name
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	return
}

func Open(driverName, dataSourceName string) (*DB, error) {
	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{DB: d}, nil
}
