package comments

import (
	"database/sql"
	"fmt"
	. "github.com/dujigui/blog/services/db"
	. "github.com/dujigui/blog/utils"
	"log"
	"time"
)

//noinspection SqlNoDataSourceInspection
const (
	tableStmt = `create table if not exists %s
(
    id      int primary key auto_increment not null,
    post_id int                            not null,
    user_id int                            not null,
    content varchar(255)                   not null,
    inspect boolean                        not null default false,
    created datetime                       not null default current_timestamp,
    updated datetime                       not null default current_timestamp on update current_timestamp
);`
	tableName = "comments"
)

func init() {
	if _, err := DB().Exec(fmt.Sprintf(tableStmt, tableName)); err != nil {
		log.Fatal("mysql ", "初始化 Table 失败 ", Params{"table": tableName}.Err(err))
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
	ByPost(pid int) ([]Comment, error)
}

type Comment struct {
	ID      int
	PostID  int
	UserID  int
	Content string
	Inspect bool
	Created time.Time
	Updated time.Time
	User    CommentUser
}

type CommentUser struct {
	ID       int
	Avatar   string
	Nickname string
}

type mysql struct {
}

func (m *mysql) Create(params Params) (int, error) {
	return Create(tableName, params)
}

func (m *mysql) Retrieve(id int) (Comment, error) {
	var c Comment
	err := Retrieve(tableName, id, func(rows *sql.Rows) error {
		return scan(&c, rows)
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
		err := scan(&c, rows)
		if err == nil {
			cs = append(cs, c)
		}
		return err
	}, limit, (page-1)*limit)
	return cs, t, err
}

func (m *mysql) ByPost(pid int) (cs []Comment, err error) {
	var c Comment
	err = Condition(tableName, "where post_id=?", func(rows *sql.Rows) error {
		err := scan(&c, rows)
		if err == nil {
			cs = append(cs, c)
		}
		return err
	}, pid)
	return
}

func scan(c *Comment, rows *sql.Rows) error {
	return rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.Inspect, &c.Created, &c.Updated)
}
