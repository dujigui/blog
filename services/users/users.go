package users

import (
	"database/sql"
	"fmt"
	. "github.com/dujigui/blog/services/db"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

//noinspection SqlNoDataSourceInspection
const (
	tableStmt = `create table if not exists %s
(
    id       int primary key auto_increment not null,
    username varchar(255),
    password varchar(255),
    admin    boolean                        not null default false,
    type     int                            not null default 2 comment '1 账号密码注册; 2 qq注册;',
    qq_id    int                            not null default 0,
    avatar   varchar(255)                   not null default '/images/avatar.svg',
    nickname varchar(255)                   not null default '',
    created  datetime                       not null default current_timestamp,
    updated  datetime                       not null default current_timestamp on update current_timestamp
);
`
	tableName  = "users"
	ViaAccount = 1
	ViaQQ      = 2
)

func init() {
	if _, err := DB().Exec(fmt.Sprintf(tableStmt, tableName)); err != nil {
		log.Fatal("mysql ", "初始化 Table 失败 ", Params{"table": tableName}.Err(err))
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
	ByUsername(username string) (User, error)
	ByQQID(id int) (User, error)
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		Logger().Fatal("users", "无法 hash 密码", Params{"password": password})
	}
	return string(hash)
}

func ComparePassword(account, password string) (User, bool) {
	u, err := UserTable().ByUsername(account)
	if err != nil {
		Logger().Trace("users", "无此用户", Params{"account": account, "password": password}.Err(err))
		return u, false
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		Logger().Trace("users", "密码不正确", Params{"account": account, "password": password}.Err(err))
		return u, false
	}

	return u, true
}

type User struct {
	ID       int
	Username string
	Password string
	Admin    bool
	Type     int
	QQID     int
	Avatar   string
	Nickname string
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
		return scan(&user, rows)
	})
	return user, err
}

func (u *mysql) Update(id int, params Params) error {
	return Update(tableName, id, params)
}

func (u *mysql) Delete(id int) error {
	return Delete(tableName, id)
}

func (u *mysql) ByUsername(username string) (user User, err error) {
	err = Condition(tableName, "where username=?", func(rows *sql.Rows) error {
		return scan(&user, rows)
	}, username)
	return
}

func (u *mysql) ByQQID(id int) (user User, err error) {
	err = Condition(tableName, "where qq_id=?", func(rows *sql.Rows) error {
		return scan(&user, rows)
	}, id)
	return
}

func scan(u *User, rows *sql.Rows) error {
	return rows.Scan(&u.ID, &u.Username, &u.Password, &u.Admin, &u.Type, &u.QQID, &u.Avatar, &u.Nickname, &u.Created, &u.Updated)
}
