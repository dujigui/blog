package users

import (
	"database/sql"
	"fmt"
	. "github.com/dujigui/blog/services/db"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"time"
)

//noinspection SqlNoDataSourceInspection
const (
	tableStmt = `create table if not exists %s
(
    id       int primary key auto_increment not null,
    username varchar(255)                   not null,
    password varchar(255)                   not null,
    created  datetime                       not null,
    updated  datetime                       not null
);`
	tableName = "users"
)

func init() {
	if _, err := DB().Exec(fmt.Sprintf(tableStmt, tableName)); err != nil {
		Logger().Fatal("mysql", "初始化 Table 失败", Params{"table": tableName, "err": err})
	}
}

func UserTable() Users {
	return &mysql{}
}

type Users interface {
	Create(params Params) (int, error)
	Retrieve(id int) (User, error)
	Update(id int, params Params) error
	Delete(id int) error
}

type User struct {
	ID       int
	Username string
	Password string
	Created  time.Time
	Updated  time.Time
}

type mysql struct {
}

func (u *mysql) Create(params Params) (int, error) {
	return Create(tableName, params)
}

func (u *mysql) Retrieve(id int) (User, error) {
	var user User
	err := Retrieve(tableName, id, func(rows *sql.Rows) error {
		return rows.Scan(&user.ID, &user.Username, &user.Password, &user.Created, &user.Updated)
	})
	return user, err
}

func (u *mysql) Update(id int, params Params) error {
	return Update(tableName, id, params)
}

func (u *mysql) Delete(id int) error {
	return Delete(tableName, id)
}
