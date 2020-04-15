package db

import (
	"database/sql"
	"errors"
	. "github.com/dujigui/blog/utils"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var (
	db             *sql.DB
	ErrIDNotExists = errors.New("无此 ID 数据")
)

func DB() *sql.DB {
	return db
}

func init() {
	var err error
	dsn := os.Getenv("BLOG_DSN")
	if dsn == "" {
		log.Fatal("dsn 为空，请设置 BLOG_DSN 环境变量")
	}
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("mysql ", "初始化数据库失败 ", Params{"dsn": dsn}.Err(err))
	}
}
