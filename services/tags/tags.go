package tags

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
    id int primary key auto_increment not null,
    name varchar(255) not null ,
    created datetime not null default current_timestamp,
    updated datetime not null default current_timestamp on update current_timestamp
);`
	tableName = "tags"
)

func init() {
	if _, err := DB().Exec(fmt.Sprintf(tableStmt, tableName)); err != nil {
		log.Fatal("mysql ", "初始化 Table 失败 ", Params{"table": tableName}.Err(err))
	}
}

func TagTable() Tags {
	return &mysql{}
}

type Tags interface {
	Create(params Params) (int, error)
	Retrieve(id int) (Tag, error)
	Update(id int, params Params) error
	Delete(id int) error
	RetrieveIDs(ids string) ([]Tag, error)
	All() ([]Tag, error)
}

type Tag struct {
	ID      int
	Name    string
	Created time.Time
	Updated time.Time
}

type mysql struct {
}

func (u *mysql) Create(params Params) (int, error) {
	return Create(tableName, params)
}

func (u *mysql) Retrieve(id int) (Tag, error) {
	var t Tag
	err := Retrieve(tableName, id, func(rows *sql.Rows) error {
		return rows.Scan(&t.ID, &t.Name, &t.Created, &t.Updated)
	})
	return t, err
}

func (u *mysql) Update(id int, params Params) error {
	return Update(tableName, id, params)
}

func (u *mysql) Delete(id int) error {
	return Delete(tableName, id)
}

// id=5,10,11
func (u *mysql) RetrieveIDs(ids string) (tags []Tag, err error) {
	var t Tag
	err = Condition(tableName, fmt.Sprintf("where id in (%s)", ids), func(rows *sql.Rows) error {
		err := rows.Scan(&t.ID, &t.Name, &t.Created, &t.Updated)
		tags = append(tags, t)
		return err
	})
	return
}

func (u *mysql) All() (tags []Tag, err error) {
	var t Tag
	err = Condition(tableName, "", func(rows *sql.Rows) error {
		err := rows.Scan(&t.ID, &t.Name, &t.Created, &t.Updated)
		tags = append(tags, t)
		return err
	})
	return
}
