package posts

import (
	"database/sql"
	"fmt"
	. "github.com/dujigui/blog/services/db"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/services/tags"
	. "github.com/dujigui/blog/utils"
	"time"
)

//noinspection SqlNoDataSourceInspection
const (
	tableStmt = `create table if not exists %s
(
    id          int primary key auto_increment not null,
    title       varchar(255),
    description varchar(255),
    cover       varchar(255),
    publish     boolean                        not null default false comment '是否发布',
    type        int                            not null,
    content     blob                           not null,
    tag_ids     varchar(255)                   not null default '',
    created     datetime                       not null default current_timestamp,
    updated     datetime                       not null default current_timestamp on update current_timestamp
);`
	tableName = "posts"
	Article   = 1
	Snippet   = 2
	Moment    = 3
)

func init() {
	if _, err := DB().Exec(fmt.Sprintf(tableStmt, tableName)); err != nil {
		Logger().Fatal("mysql", "初始化 Table 失败", Params{"table": tableName, "err": err})
	}
}

func PostTable() Posts {
	return &mysql{}
}

type Posts interface {
	Create(params Params) (int, error)
	Retrieve(id int) (Post, error)
	Update(id int, params Params) error
	Delete(id int) error
	Page(page, limit int, all bool) ([]Post, int, error)
	Count() (int, error)
	Latest() (Post, error)
	All() ([]Post, error)
}

//todo 引进author的概念
type Post struct {
	ID          int
	Title       string
	Description string
	Cover       string
	Content     []byte
	Publish     bool
	Type        int
	TagIDs      string // 1,2,3
	Tags        []Tag
	Created     time.Time
	Updated     time.Time
}

type mysql struct {
}

func (m *mysql) Create(params Params) (int, error) {
	return Create(tableName, params)
}

func (m *mysql) Retrieve(id int) (Post, error) {
	var p Post
	err := Retrieve(tableName, id, func(rows *sql.Rows) error {
		return scan(&p, rows)
	})
	return p, err
}

func (m *mysql) Update(id int, params Params) error {
	return Update(tableName, id, params)
}

func (m *mysql) Delete(id int) error {
	return Delete(tableName, id)
}

func (m *mysql) Page(page, limit int, all bool) ([]Post, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit >= 50 {
		limit = 50
	}

	var posts = make([]Post, 0)
	var p Post

	var where string
	if all {
		where = "order by created desc"
	} else {
		where = "where publish=true order by created desc"
	}
	t, err := Page(tableName, where, func(rows *sql.Rows) error {
		err := scan(&p, rows)
		if err == nil {
			posts = append(posts, p)
		}
		return err
	}, limit, (page-1)*limit)
	return posts, t, err
}

func (m *mysql) Count() (int, error) {
	return Count(tableName, "")
}

func (m *mysql) Latest() (p Post, err error) {
	err = Condition(tableName, "order by updated desc", func(rows *sql.Rows) error {
		return scan(&p, rows)
	})
	return
}

func (m *mysql) All() ([]Post, error) {
	var posts = make([]Post, 0)
	var p Post
	err := Condition(tableName, "", func(rows *sql.Rows) error {
		err := scan(&p, rows)
		if err == nil {
			posts = append(posts, p)
		}
		return err
	})
	return posts, err
}

func scan(p *Post, rows *sql.Rows) error {
	return rows.Scan(&p.ID, &p.Title, &p.Description, &p.Cover, &p.Publish, &p.Type, &p.Content, &p.TagIDs, &p.Created, &p.Updated)
}
