package posts

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
    id           int primary key auto_increment not null,
    title        varchar(255),
    description  varchar(255),
    cover        varchar(255),
    created      datetime                       not null,
    updated      datetime                       not null,
    is_published boolean                        not null,
    type         int                            not null,
    content      blob                           not null
);`
	tableName = "posts"
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
	Page(page int) ([]Post, int, error)
}

type Post struct {
	ID          int
	Title       string
	Description string
	Cover       string
	Content     []byte
	Created     time.Time
	Updated     time.Time
	IsPublished bool
	Type        int
}

type mysql struct {
}

func (m *mysql) Create(params Params) (int, error) {
	return Create(tableName, params)
}

func (m *mysql) Retrieve(id int) (Post, error) {
	var p Post
	err := Retrieve(tableName, id, func(rows *sql.Rows) error {
		return rows.Scan(&p.ID, &p.Title, &p.Description, &p.Cover, &p.Created, &p.Updated, &p.IsPublished, &p.Type, &p.Content)
	})
	return p, err
}

func (m *mysql) Update(id int, params Params) error {
	return Update(tableName, id, params)
}

func (m *mysql) Delete(id int) error {
	return Delete(tableName, id)
}

func (m *mysql) Page(page int) ([]Post, int, error) {
	var posts = make([]Post, 0)
	var p Post
	t, err := Page(tableName, "", func(rows *sql.Rows) error {
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Cover, &p.Created, &p.Updated, &p.IsPublished, &p.Type, &p.Content)
		if err == nil {
			posts = append(posts, p)
		}
		return err
	}, 10, (page-1)*10)
	return posts, t, err
}
