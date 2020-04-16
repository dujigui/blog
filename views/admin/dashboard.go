package admin

import (
	"github.com/dujigui/blog/gateway"
	. "github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/comments"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/services/posts"
	. "github.com/dujigui/blog/services/users"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"time"
)

type adminCtrl struct {
	Ctx iris.Context
}

// GET /admin
func (c *adminCtrl) Get() mvc.View {
	pc, err := PostTable().Count()
	if err != nil {
		Logger().Error("adminCtrl", "获取文章总数失败", Params{"err": err})
	}
	cc, err := CommentTable().Count()
	if err != nil {
		Logger().Error("adminCtrl", "获取评论总数失败", Params{"err": err})
	}

	var dp int
	lp, err := PostTable().Latest()
	if err == nil && lp.ID > 0 {
		dp = int(time.Now().Sub(lp.Updated).Hours() / 24)
	} else {
		dp = 0
	}
	do := int(time.Now().Sub(time.Unix(Pref().Init, 0)).Hours() / 24)

	return mvc.View{
		Name: "admin/html/dashboard.html",
		Data: iris.Map{
			"tab":           dashboard,
			"PostNumber":    pc,
			"CommentNumber": cc,
			"DaysLastPost":  dp,
			"DayOnline":     do,
		},
	}
}

func Info(ctx iris.Context) {
	ok, uid, admin := gateway.Info(ctx)

	if !ok || !admin || uid <= 0 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(Result(false, "StatusUnauthorized", nil))
		Logger().Warning("adminCtrl", "未授权访问用户信息", nil)
		return
	}

	u, err := UserTable().Retrieve(uid)
	if err != nil || u.ID <= 0 {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "StatusInternalServerError", nil))
		Logger().Warning("adminCtrl", "无此用户信息", Params{"id": uid})
		return
	}

	type formUser struct {
		ID       int    `json:"id"`
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
	}
	ctx.JSON(Result(true, "ok", formUser{
		ID:       u.ID,
		Avatar:   u.Avatar,
		Nickname: u.Nickname,
	}))
}
