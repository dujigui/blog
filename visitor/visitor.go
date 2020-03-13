package visitor

import "github.com/kataras/iris/v12"

func Visitor(app *iris.Application) {
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>index page</h1>")
	})
}
