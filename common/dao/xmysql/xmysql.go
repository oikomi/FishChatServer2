package xmysql

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/conf"
	"golang.org/x/net/context"
	"time"
)

var (
	ErrStmtNil = errors.New("prepare failed and stmt nil")
)

type DB struct {
	*sql.DB
	env string
}

type Tx struct {
	tx *sql.Tx
}

type Stmt struct {
	stmt  *sql.Stmt
	query string
	tx    bool
	env   string
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

func (db *DB) Exec(c context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.DB.Exec(query, args...)
}

func (db *DB) Ping(c context.Context) error {
	return db.DB.Ping()
}

func (db *DB) Prepare(query string) (*Stmt, error) {
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt, query: query, env: db.env}, nil
}

func (db *DB) Prepared(query string) (stmt *Stmt) {
	stmt = &Stmt{query: query, env: db.env}
	s, err := db.DB.Prepare(query)
	if err == nil {
		stmt.stmt = s
		return
	}
	go func() {
		for {
			s, err := db.DB.Prepare(query)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			stmt.stmt = s
			return
		}
	}()
	return
}

func (db *DB) Query(c context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.DB.Query(query, args...)
}

func (db *DB) QueryRow(c context.Context, query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(query, args...)
}

func (s *Stmt) Close() error {
	if s.stmt != nil {
		return s.stmt.Close()
	}
	return nil
}

func (s *Stmt) Exec(c context.Context, args ...interface{}) (sql.Result, error) {
	if s.stmt == nil {
		return nil, ErrStmtNil
	}
	return s.stmt.Exec(args...)
}

func (s *Stmt) Query(c context.Context, args ...interface{}) (*sql.Rows, error) {
	if s.stmt == nil {
		return nil, ErrStmtNil
	}
	return s.stmt.Query(args...)
}

func (s *Stmt) QueryRow(c context.Context, args ...interface{}) (*sql.Row, error) {
	if s.stmt == nil {
		return nil, ErrStmtNil
	}
	return s.stmt.QueryRow(args...), nil
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.tx.Exec(query, args...)
}

func (tx *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.tx.Query(query, args...)
}

func (tx *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	return tx.tx.QueryRow(query, args...)
}

func (tx *Tx) Stmt(stmt *Stmt) *Stmt {
	st := tx.tx.Stmt(stmt.stmt)
	return &Stmt{stmt: st, query: stmt.query, tx: true}
}

func (tx *Tx) Prepare(query string) (*Stmt, error) {
	stmt, err := tx.tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &Stmt{stmt: stmt, query: query, tx: true}, nil
}
