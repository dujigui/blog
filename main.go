package main

import (
	"fmt"
	. "github.com/dujigui/blog/admin"
	. "github.com/dujigui/blog/gateway"
	. "github.com/dujigui/blog/services/cfg"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/visitor"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(ReqLogger())
	app.Use(Gateway)
	Admin(app)
	Visitor(app)
	host := Config().GetString("host")
	port := Config().GetString("port")
	_ = app.Run(iris.Addr(fmt.Sprintf("%s:%s", host, port)), iris.WithoutServerError(iris.ErrServerClosed))
}
