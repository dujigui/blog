package files

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
    original varchar(255)                   null,
    hashed   varchar(255)                   null,
    created datetime not null default current_timestamp,
    updated datetime not null default current_timestamp on update current_timestamp
);`
	tableName = "files"
)

func init() {
	if _, err := DB().Exec(fmt.Sprintf(tableStmt, tableName)); err != nil {
		Logger().Fatal("mysql", "初始化 Table 失败", Params{"table": tableName, "err": err})
	}
}

func FileTable() Files {
	return &mysql{}
}

type Files interface {
	Create(params Params) (int, error)
	Retrieve(id int) (File, error)
	Update(id int, params Params) error
	Delete(id int) error
	Page(page, limit int) ([]File, int, error)
}

type File struct {
	ID       int
	Original string
	Hashed   string
	Created  time.Time
	Updated  time.Time
}

type mysql struct {
}

func (m *mysql) Create(params Params) (int, error) {
	return Create(tableName, params)
}

func (m *mysql) Retrieve(id int) (File, error) {
	var f File
	err := Retrieve(tableName, id, func(rows *sql.Rows) error {
		return rows.Scan(&f.ID, &f.Original, &f.Hashed, &f.Created, &f.Updated)
	})
	return f, err
}

func (m *mysql) Update(id int, params Params) error {
	return Update(tableName, id, params)
}

func (m *mysql) Delete(id int) error {
	return Delete(tableName, id)
}

func (m *mysql) Page(page, limit int) ([]File, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit >= 50 {
		limit = 50
	}

	var files = make([]File, 0)
	var f File
	t, err := Page(tableName, "order by created desc", func(rows *sql.Rows) error {
		err := rows.Scan(&f.ID, &f.Original, &f.Hashed, &f.Created, &f.Updated)
		if err == nil {
			files = append(files, f)
		}
		return err
	}, limit, (page-1)*limit)
	return files, t, err
}
