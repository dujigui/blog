package visitor

import (

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

const (
	home = "home"
	search = "search"
	about = "about"
)

func Visitor(app *iris.Application) {
	app.HandleDir("/visitor/css", "views/web/visitor/css")
	app.HandleDir("/visitor/js", "views/web/visitor/js")

	hp := app.Party("/").Layout("visitor/visitor.html")
	mvc.New(hp).Handle(new(visitorCtr))
}