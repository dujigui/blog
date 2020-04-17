package main

import (
	. "github.com/dujigui/blog/gateway"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	. "github.com/dujigui/blog/views"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"os"
	"strings"
	"time"
)

func main() {
	if strings.EqualFold(os.Getenv("BLOG_DEBUG"), "true") {
		Logger().Trace("main", "使用 BLOG_DEBUG 模式运行", nil)
	}
	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(ReqLogger())
	app.Use(Gateway)
	Views(app)
	if err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed)); err != nil {
		Logger().Error("main", "运行终止", Params{"at": time.Now()}.Err(err))
	}
}
