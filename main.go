package main

import (
	. "github.com/dujigui/blog/gateway"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/views"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(ReqLogger())
	app.Use(Gateway)
	Views(app)
	_ = app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

/*func renameFile() {
	rows, err := db.DB().Query("select id,original,hashed from files")
	if err != nil {
		log.Fatal(err)
	}
	var id int
	var fn, hashed string
	for rows.Next() {
		rows.Scan(&id, &fn, &hashed)
		db.DB().Exec("update files set hashed=? where id=?", hashed + filepath.Ext(fn), id)
	}
	rows.Close()
}

func renameCover() {
	rows, err := db.DB().Query("select id,cover from posts")
	if err != nil {
		log.Fatal(err)
	}
	var id int
	var cover string
	for rows.Next() {
		rows.Scan(&id, &cover)
		if strings.HasPrefix(cover, "/attachment") {
			cover = strings.Replace(cover, "/attachment", "/files", -1)
			db.DB().Exec("update posts set cover=? where id=?",cover, id)
		}
	}
	rows.Close()
}

func reType() {
	rows, err := db.DB().Query("select id,title, description, type from posts")
	if err != nil {
		log.Fatal(err)
	}
	var id, t int
	var title, description string
	for rows.Next() {
		rows.Scan(&id, &title, &description, &t)
		if title == "" {
			if description == "" {
				//fmt.Printf("should=3 now=%d title=%s desc=%s\n", t, title, description)
				db.DB().Exec("update posts set type=? where id=?",3, id)
				continue
			}
			//fmt.Printf("should=2 now=%d title=%s desc=%s\n", t, title, description)
			db.DB().Exec("update posts set type=? where id=?",2, id)
			continue
		}
		//fmt.Printf("should=1 now=%d title=%s desc=%s\n", t, title, description)
		db.DB().Exec("update posts set type=? where id=?",1, id)
	}
	rows.Close()
}*/
