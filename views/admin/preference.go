package admin

import (
	. "github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type prefCtrl struct {
	Ctx iris.Context
}

// GET /admin/preferences
func (c *prefCtrl) Get() mvc.View {
	return mvc.View{
		Name: "admin/html/preference.html",
		Data: iris.Map{
			"tab": preference,
		},
	}
}

// GET /admin/preferences/get
func getPref(ctx iris.Context) {
	ctx.JSON(Result(true, "ok", *Pref()))
}

// PATCH /admin/preferences
func patchPref(ctx iris.Context) {
	p := &Preferences{}
	if err := ctx.ReadJSON(p); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "无法读取body", nil))
		Logger().Error("prefCtrl", "无法读取body", Params{}.Err(err))
		return
	}

	pp := Pref()
	pp.BlogName = p.BlogName
	pp.AdminPageName = p.AdminPageName
	pp.AboutPostID = p.AboutPostID
	pp.Email = p.Email

	if err := pp.Save(); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "更新设置失败", nil))
		Logger().Error("tagCtrl", "更新设置失败", Params{"p": p}.Err(err))
		return
	}

	ctx.JSON(Result(true, "ok", nil))
}
