package visitor

import (
	"fmt"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/services/users"
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

// todo 检查所有用户的输入，防止sql注入和xss攻击
func Login(ctx iris.Context) {
	type form struct {
		Account  string `form:"account"`
		Password string `form:"password"`
	}
	var f form
	if err := ctx.ReadJSON(&f); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(true, "无法读取 body", nil))
		return
	}

	if f.Account == "" || f.Password == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(true, "无法读取账号/密码", nil))
		return
	}

	u, ok := ComparePassword(f.Account, f.Password)
	if !ok {
		fmt.Printf("account=%s password=%s\n", f.Account, f.Password)
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(true, "账号/密码不正确", nil))
		return
	}

	ctx.SetCookieKV("token", CreateToken(u.ID, u.Admin, CookieExpire), iris.CookieExpires(CookieExpire))
	ctx.JSON(Result(true, "ok", nil))
}

func Logout(ctx iris.Context) {
	ctx.RemoveCookie("token")
	ctx.JSON(Result(true, "ok", nil))
}
