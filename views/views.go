package views

import (
	"github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/views/admin"
	. "github.com/dujigui/blog/views/visitor"
	"github.com/kataras/iris/v12"
	"os"
	"strings"
	"time"
)

func Views(app *iris.Application) {
	app.HandleDir("/files", "data/files")
	app.HandleDir("/", "data/favicon/")
	app.HandleDir("/layui", "views/layui")
	app.HandleDir("/prism", "views/prism")
	app.HandleDir("/images", "views/images")


	tplEngine := iris.HTML("views/web", ".html")
	isDebug := strings.EqualFold(os.Getenv("BLOG_DEBUG"), "true")
	tplEngine.Reload(isDebug)
	{
		//公共部分
		tplEngine.AddFunc("pref", services.Pref)
	}
	{
		// admin 管理前端部分
	}
	{
		// visitor 用户前端部分
		tplEngine.AddFunc("date", date)
		tplEngine.AddFunc("add", add)
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