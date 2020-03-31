package visitor

import (

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
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

	hp := app.Party("/").Layout("visitor/visitor.html")
	mvc.New(hp).Handle(new(homeCtrl))

	sp := hp.Party("/search")
	mvc.New(sp).Handle(new(searchCtrl))

	dp := hp.Party("/posts")
	mvc.New(dp).Handle(new(detailCtrl))

	ap := hp.Party("/about")
	mvc.New(ap).Handle(new(aboutCtrl))
}