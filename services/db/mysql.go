package db

import (
	"database/sql"
	"fmt"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"strings"
)

type Scanner func(rows *sql.Rows) error

func Create(table string, params Params) (int, error) {
	var keys []string
	var values []interface{}
	var qm []string
	for k, v := range params {
		keys = append(keys, k)
		values = append(values, v)
		qm = append(qm, "?")
	}

	tpl := "insert into %s(%s) values(%s)"
	q := fmt.Sprintf(tpl, table, strings.Join(keys, ","), strings.Join(qm, ","))

	Logger().Trace("mysql", "准备插入数据", Params{"table": table, "q": q}.AddAll(params))

	stmt, err := DB().Prepare(q)
	if err != nil {
		Logger().Error("mysql", "准备插入数据失败", Params{"table": table, "err": err}.AddAll(params))
		return 0, err
	}

	r, err := stmt.Exec(values...)
	if err != nil {
		Logger().Error("mysql", "插入数据失败", Params{"table": table, "err": err}.AddAll(params))
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		Logger().Error("mysql", "获取插入数据结果失败", Params{"table": table, "err": err}.AddAll(params))
		return 0, err
	}
	Logger().Trace("mysql", "插入数据成功", Params{"table": table}.AddAll(params))
	return int(id), nil
}

func Retrieve(table string, id int, s Scanner) error {
	rows, err := DB().Query(fmt.Sprintf("select * from %s where id=?", table), id)
	if err != nil {
		Logger().Error("mysql", "获取数据失败", Params{"table": table, "id": id, "err": err})
		return err
	}

	if rows.Next() {
		err = s(rows)
		if err == nil {
			Logger().Trace("mysql", "扫描数据成功", Params{"table": table, "id": id})
		} else {
			Logger().Trace("mysql", "扫描数据失败", Params{"table": table, "id": id, "err": err})
		}
		rows.Close()
		return err
	}
	rows.Close()
	Logger().Error("mysql", "无法获取不存在的数据", Params{"table": table, "id": id, "err": err})
	return ErrIDNotExists
}

func Update(table string, id int, params Params) error {
	var keys []string
	var values []interface{}
	for k, v := range params {
		keys = append(keys, fmt.Sprintf("%s=?", k))
		values = append(values, v)
	}
	values = append(values, id)

	tpl := "update %s set %s where id=?"
	q := fmt.Sprintf(tpl, table, strings.Join(keys, ","))

	Logger().Trace("mysql", "准备更新数据", Params{"table": table, "q": q}.AddAll(params))

	stmt, err := DB().Prepare(q)
	if err != nil {
		Logger().Error("mysql", "准备更新数据失败", Params{"table": table, "err": err}.AddAll(params))
		return err
	}

	r, err := stmt.Exec(values...)
	if err != nil {
		Logger().Error("mysql", "更新数据失败", Params{"table": table, "err": err}.AddAll(params))
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		Logger().Error("mysql", "获取更新数据结果失败", Params{"table": table, "id": id, "err": err}.AddAll(params))
		return err
	}
	if n != 1 {
		Logger().Error("mysql", "更新大于 1 行的数据", Params{"table": table, "id": id, "n": n}.AddAll(params))
	}
	Logger().Trace("mysql", "更新数据成功", Params{"table": table, "id": id}.AddAll(params))
	return nil
}

func Delete(table string, id int) error {
	r, err := DB().Exec(fmt.Sprintf("delete from %s where id=?", table), id)
	if err != nil {
		Logger().Error("mysql", "删除数据失败", Params{"table": table, "id": id, "err": err})
		return err
	}
	n, err := r.RowsAffected()
	if err != nil {
		Logger().Error("mysql", "获取删除数据结果失败", Params{"table": table, "id": id, "err": err})
		return err
	}
	if n != 1 {
		Logger().Error("mysql", "删除大于 1 行的数据", Params{"table": table, "id": id, "n": n})
	}
	Logger().Trace("mysql", "删除数据成功", Params{"table": table, "id": id})
	return nil
}

func Condition(table, condition string, s Scanner, values ...interface{}) error {
	if condition == "" {
		condition = "true"
	}
	q := fmt.Sprintf("select * from %s where %s", table, condition)
	fmt.Println(q)
	rows, err := DB().Query(q, values...)
	if err != nil {
		Logger().Error("mysql", "准备获取数据列表失败", Params{"table": table, "values": values, "err": err})
		return err
	}

	for rows.Next() {
		err = s(rows)
		if err != nil {
			Logger().Error("mysql", "获取数据列表失败", Params{"table": table, "values": values, "err": err})
			rows.Close()
			return err
		}
	}
	rows.Close()
	Logger().Trace("mysql", "获取数据列表成功", Params{"table": table, "values": values})
	return nil
}

func Page(table, condition string, s Scanner, limit, offset int, values ...interface{}) (int, error) {
	if condition == "" {
		condition = "true"
	}

	var total int
	rows, err := DB().Query(fmt.Sprintf("select count(*) from %s where %s", table, condition), values...)
	if err != nil {
		Logger().Error("mysql", "获取数据总数失败", Params{"table": table, "condition": condition, "values": values, "limit": limit, "offset": offset, "err": err})
		return 0, err
	}
	if rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			Logger().Error("mysql", "扫描数据总数失败", Params{"table": table, "condition": condition, "values": values, "limit": limit, "offset": offset, "err": err})
			rows.Close()
			return 0, err
		}
	}
	rows.Close()

	values = append(values, limit)
	values = append(values, offset)
	rows, err = DB().Query(fmt.Sprintf("select * from %s where %s limit ? offset ?", table, condition), values...)
	if err != nil {
		Logger().Error("mysql", "准备获取分页列表失败", Params{"table": table, "values": values, "err": err})
		return total, err
	}

	for rows.Next() {
		err = s(rows)
		if err != nil {
			Logger().Error("mysql", "获取分页列表失败", Params{"table": table, "values": values, "err": err})
			return total, err
		}
	}
	rows.Close()
	Logger().Trace("mysql", "获取分页列表成功", Params{"table": table, "values": values})
	return total, nil
}
