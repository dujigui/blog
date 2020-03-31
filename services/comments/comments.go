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
    created datetime not null default current_timestamp,
    updated datetime not null default current_timestamp on update current_timestamp
)`
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
	Count() (int, error)
	Page(page, limit int) ([]Comment, int, error)
	RetrieveByPost(pid int) ([]Comment, error)
}

// todo 添加评论审核功能
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

func (m *mysql) Count() (int, error) {
	return Count(tableName, "")
}

func (m *mysql) Page(page, limit int) ([]Comment, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit >= 50 {
		limit = 50
	}

	var cs = make([]Comment, 0)
	var c Comment
	t, err := Page(tableName, "order by created desc", func(rows *sql.Rows) error {
		err := rows.Scan(&c.ID, &c.PostID, &c.Content, &c.Created, &c.Updated)
		if err == nil {
			cs = append(cs, c)
		}
		return err
	}, limit, (page-1)*limit)
	return cs, t, err
}

func (m *mysql) RetrieveByPost(pid int) (cs []Comment, err error) {
	var c Comment
	err = Condition(tableName, "where post_id=?", func(rows *sql.Rows) error {
		err := rows.Scan(&c.ID, &c.PostID, &c.Content, &c.Created, &c.Updated)
		if err == nil {
			cs = append(cs, c)
		}
		return err
	}, pid)
	return
}
