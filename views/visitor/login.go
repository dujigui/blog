package visitor

import (
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"time"
)

type loginCtrl struct {
	Ctx iris.Context
}

func (c *loginCtrl) Get() mvc.View {
	re := c.Ctx.URLParam("redirect")
	if re == "" {
		re = "/"
	}

	state := QQState{
		Timestamp: time.Now().Unix(),
		Redirect:  re,
	}
	d, err := EncodeState(state)
	if err != nil {
		Logger().Error("login", "无法创建 QQ 登录状态", Params{}.Err(err))
		return ErrMsg("无法创建 QQ 登录状态")
	}

	return mvc.View{
		Name: "visitor/html/login.html",
		Data: iris.Map{
			"qqState": d,
		},
	}
}
