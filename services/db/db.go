package db

import (
	"database/sql"
	"errors"
	"github.com/dujigui/blog/services/cfg"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
	ErrIDNotExists = errors.New("无此 ID 数据")
)

func DB() *sql.DB {
	return db
}

func init() {
	var err error
	dsn := cfg.Config().GetString("dsn")
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		Logger().Fatal("mysql", "初始化数据库失败", Params{"dsn": dsn, "err": err})
	}
}
