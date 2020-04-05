package visitor

import (

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/view"
)

const (
	home = "home"
	search = "search"
	about = "about"
	detail = "detail"
	none = "none"
)

func Visitor(app *iris.Application) {
	app.HandleDir("/visitor/css", "views/web/visitor/css")
	app.HandleDir("/visitor/js", "views/web/visitor/js")

	ip := app.Party("/init").Layout(view.NoLayout)
	mvc.New(ip).Handle(new(initCtrl))
	ip.Post("/", Init)

	hp := app.Party("/").Layout("visitor/visitor.html")
	mvc.New(hp).Handle(new(homeCtrl))

	sp := hp.Party("/search")
	mvc.New(sp).Handle(new(searchCtrl))

	dp := hp.Party("/posts")
	mvc.New(dp).Handle(new(detailCtrl))

	ap := hp.Party("/about")
	mvc.New(ap).Handle(new(aboutCtrl))

	lp := app.Party("/login").Layout(view.NoLayout)
	mvc.New(lp).Handle(new(loginCtrl))
	app.Post("/login", Login)
	app.Delete("/logout", Logout)

	app.Get("/qq", qq)
	app.Post("/comments", postComment)
}