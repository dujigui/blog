package admin

import (
	. "github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/comments"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/services/posts"
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
	if err == nil {
		dp = int(time.Now().Sub(lp.Updated).Hours() / 24)
	}
	do := int(time.Now().Sub(time.Unix(Pref().Init, 0)).Hours() / 24)

	return mvc.View{
		Name: "admin/html/dashboard.html",
		Data: iris.Map{
			"tab":            dashboard,
			"PostNumber":     pc,
			"CommentNumber":  cc,
			"DaysLastPost":   dp,
			"DayOnline":      do,
		},
	}
}

