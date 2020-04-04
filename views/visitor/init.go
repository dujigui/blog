package visitor

import (
	"errors"
	. "github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/users"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"time"
)

type initCtrl struct {
	Ctx iris.Context
}

func (c *initCtrl) Get() mvc.View {
	return mvc.View{
		Name: "visitor/html/init.html",
		Data: iris.Map{},
	}
}

func Init(ctx iris.Context) {
	type form struct {
		BlogName      string
		AdminPageName string
		Email         string
		Salt          string
		Account       string
		Password      string
	}
	var f form
	if err := ctx.ReadJSON(&f); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("无法读取 body")
		return
	}
	if f.Salt == "" || f.Account == "" || f.Password == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("请按要求填写数据")
		return
	}

	if err := admin(f.Account, f.Password); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	Pref().BlogName = f.BlogName
	Pref().AdminPageName = f.AdminPageName
	Pref().Email = f.Email
	Pref().Salt = f.Salt
	Pref().Init = time.Now().Unix()
	Pref().Save()

	ctx.JSON(Result(true, "ok", nil))
}

func admin(account, password string) error {
	u, err := UserTable().ByUsername(account)
	if u.ID != 0 {
		return errors.New("用户名已存在")
	}
	id, err := UserTable().Create(Params{"username": account, "password": HashPassword(password), "admin": true})
	if err != nil {
		return err
	}
	if id <= 0 {
		return errors.New("创建管理员失败")
	}
	return nil
}
