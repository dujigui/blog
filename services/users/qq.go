package users

import (
	"database/sql"
	"fmt"
	. "github.com/dujigui/blog/services/db"
	. "github.com/dujigui/blog/utils"
	"log"
)

//noinspection SqlNoDataSourceInspection
const (
	tableQQStmt = `create table if not exists %s
(
    id             int primary key auto_increment not null,
    open_id        varchar(255),
    access_token   varchar(255),
    nickname       varchar(255),
    gender         varchar(50),
    gender_type    int,
    province       varchar(255),
    city           varchar(255),
    year           varchar(255),
    constellation  varchar(255),
    figureurl      varchar(255),
    figureurl_1    varchar(255),
    figureurl_2    varchar(255),
    figureurl_qq_1 varchar(255),
    figureurl_qq_2 varchar(255),
    figureurl_qq   varchar(255),
    figureurl_type varchar(10)
);`
	tableQQName = "qq"
)

func init() {
	if _, err := DB().Exec(fmt.Sprintf(tableQQStmt, tableQQName)); err != nil {
		log.Fatal("mysql ", "初始化 Table 失败 ", Params{"table": tableQQName}.Err(err))
	}
}

func QQTable() QQ {
	return &Impl{}
}

type QQ interface {
	Create(params Params) (int, error)
	Retrieve(id int) (QQInfo, error)
	Update(id int, params Params) error
	ByOpenID(openID string) (QQInfo, error)
}

type QQInfo struct {
	ID            int
	OpenID        string
	AccessToken   string

	Ret           int    `json:"ret"`
	Msg           string `json:"msg"`
	Nickname      string `json:"nickname"`
	Gender        string `json:"gender"`
	GenderType    int    `json:"gender_type"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Year          string `json:"year"`
	Constellation string `json:"constellation"`
	Figureurl     string `json:"figureurl"`
	Figureurl1    string `json:"figureurl_1"`
	Figureurl2    string `json:"figureurl_2"`
	FigureurlQq1  string `json:"figureurl_qq_1"`
	FigureurlQq2  string `json:"figureurl_qq_2"`
	FigureurlQq   string `json:"figureurl_qq"`
	FigureurlType string `json:"figureurl_type"`
}

type Impl struct {
}

func (u *Impl) Create(params Params) (int, error) {
	return Create(tableQQName, params)
}

func (u *Impl) Retrieve(id int) (qi QQInfo, err error) {
	err = Retrieve(tableQQName, id, func(rows *sql.Rows) error {
		return rows.Scan(
			&qi.ID, &qi.OpenID, &qi.AccessToken, &qi.Nickname, &qi.Gender, &qi.GenderType, &qi.Province, &qi.City, &qi.Year, &qi.Constellation,
			&qi.Figureurl, &qi.Figureurl1, &qi.Figureurl2, &qi.FigureurlQq1, &qi.FigureurlQq2, &qi.FigureurlQq, &qi.FigureurlType)
	})
	return
}

func (u *Impl) Update(id int, params Params) error {
	return Update(tableQQName, id, params)
}

func (u *Impl) ByOpenID(openID string) (qi QQInfo, err error) {
	err = Condition(tableQQName, "where open_id=?", func(rows *sql.Rows) error {
		return rows.Scan(
			&qi.ID, &qi.OpenID, &qi.AccessToken, &qi.Nickname, &qi.Gender, &qi.GenderType, &qi.Province, &qi.City, &qi.Year, &qi.Constellation,
			&qi.Figureurl, &qi.Figureurl1, &qi.Figureurl2, &qi.FigureurlQq1, &qi.FigureurlQq2, &qi.FigureurlQq, &qi.FigureurlType)
	}, openID)
	return
}
