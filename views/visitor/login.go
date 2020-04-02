package visitor

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type loginCtrl struct {
	Ctx iris.Context
}

func (c *loginCtrl) Get() mvc.View{
	return mvc.View{
		Name: "visitor/html/login.html",
		Data: iris.Map{
			//todo 加盐，存放进入的页面，保存在redis，检验，登录成功后做对应的跳转
			"qqState": "test",
		},
	}
}
