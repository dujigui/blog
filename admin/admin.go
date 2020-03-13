package admin

import "github.com/kataras/iris/v12"

func Admin(app *iris.Application) {
	admin := app.Party("/admin")
	{
		admin.Get("/", func(ctx iris.Context) {
			ctx.JSON(iris.Map{"message": "admin page"})
		})
	}
}
