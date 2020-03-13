package comments

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
    id      int primary key auto_increment not null,
    post_id int                            not null,
    content varchar(255)                   not null,
    created datetime                       not null,
    updated datetime                       not null
);`
	tableName = "comments"
)

func init() {
	if _, err := DB().Exec(fmt.Sprintf(tableStmt, tableName)); err != nil {
		Logger().Fatal("mysql", "初始化 Table 失败", Params{"table": tableName, "err": err})
	}
}

func CommentTable() Comments {
	return &mysql{}
}

type Comments interface {
	Create(params Params) (int, error)
	Retrieve(id int) (Comment, error)
	Update(id int, params Params) error
	Delete(id int) error
}

type Comment struct {
	ID      int
	PostID  int
	Content string
	Created time.Time
	Updated time.Time
}

type mysql struct {
}

func (m *mysql) Create(params Params) (int, error) {
	return Create(tableName, params)
}

func (m *mysql) Retrieve(id int) (Comment, error) {
	var c Comment
	err := Retrieve(tableName, id, func(rows *sql.Rows) error {
		return rows.Scan(&c.ID, &c.PostID, &c.Content, &c.Created, &c.Updated)
	})
	return c, err
}

func (m *mysql) Update(id int, params Params) error {
	return Update(tableName, id, params)
}

func (m *mysql) Delete(id int) error {
	return Delete(tableName, id)
}
