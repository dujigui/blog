package visitor

import (
	"github.com/dujigui/blog/services/cfg"
	"github.com/kataras/iris/v12"
)

func Visitor(app *iris.Application) {
	app.HandleDir("/", cfg.Config().GetString("favicon"))
	app.Get("/posts/{id:int}", GetBy)
}
