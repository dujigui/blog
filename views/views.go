package views

import (
	"github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/cfg"
	. "github.com/dujigui/blog/views/admin"
	. "github.com/dujigui/blog/views/visitor"
	"github.com/kataras/iris/v12"
	"time"
)

func Views(app *iris.Application) {
	app.HandleDir("/files", Config().GetString("files"))
	app.HandleDir("/", "views/favicon")
	app.HandleDir("/layui", "views/layui")
	app.HandleDir("/prism", "views/prism")
	app.HandleDir("/images", "views/images")


	tplEngine := iris.HTML("views/web", ".html")
	tplEngine.Reload(!Config().GetBool("production"))
	{
		//公共部分
		tplEngine.AddFunc("pref", services.Pref)
	}
	{
		// admin 管理前端部分
	}
	{
		// visitor 用户前端部分
		//tplEngine.AddFunc("submodule", submodule)
		//tplEngine.AddFunc("exist", exists)
		//tplEngine.AddFunc("string", str)
		//tplEngine.AddFunc("truncate", truncate)
		//tplEngine.AddFunc("truncateContent", content)
		tplEngine.AddFunc("date", date)
		tplEngine.AddFunc("add", add)
		//tplEngine.AddFunc("tagSelected", tagSelected)
	}
	app.RegisterView(tplEngine)

	Admin(app)
	Visitor(app)
}

func date(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func add(i1, i2 int) int {
	return i1+i2
}