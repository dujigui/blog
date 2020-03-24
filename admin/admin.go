package admin

import (
	. "github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/cfg"
	. "github.com/dujigui/blog/services/comments"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/services/posts"
	. "github.com/dujigui/blog/services/tags"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"time"
)

const (
	dashboard  = "dashboard"
	postList   = "post_list"
	postEditor = "post_editor"
	file       = "file"
	tag        = "tag"
	comment    = "comment"
	preference = "preference"
)

func Admin(app *iris.Application) {
	pages(app)
}

func pages(app *iris.Application) {
	tplEngine := iris.HTML("admin/html", ".html")
	tplEngine.Reload(!Config().GetBool("production"))
	tplEngine.AddFunc("getConfig", Pref)
	tplEngine.AddFunc("submodule", submodule)
	tplEngine.AddFunc("exist", exists)
	tplEngine.AddFunc("string", str)
	tplEngine.AddFunc("truncate", truncate)
	tplEngine.AddFunc("truncateContent", content)
	tplEngine.AddFunc("date", date)
	tplEngine.AddFunc("tagSelected", tagSelected)
	app.RegisterView(tplEngine)
	app.HandleDir("/backyard", "admin/html/backyard")
	app.HandleDir("/layui", "admin/html/layui")

	pa := app.Party("/admin").Layout("layout/admin.html")
	mvc.New(pa).Handle(new(adminCtrl))

	pp := pa.Party("/posts")
	mvc.New(pp).Handle(new(postsCtrl))
	pp.Get("/list", listPost)
	pp.Delete("/{id:int}", delPost)
	pp.Get("/{id:int}", getPost)

	fp := pa.Party("/files")
	mvc.New(fp).Handle(new(filesCtrl))
	fp.Post("/", iris.LimitRequestBodySize(10<<20), uploadFile)
	fp.Get("/list", listFile)
	fp.Delete("/{id:int}", delFile)

	tp := pa.Party("/tags")
	mvc.New(tp).Handle(new(tagCtrl))
	tp.Get("/list", listTag)
	tp.Post("/", createTag)
	tp.Delete("/{id:int}", delTag)
	tp.Patch("/{id:int}", patchTag)

	cp := pa.Party("/comments")
	mvc.New(cp).Handle(new(cmtCtrl))
	cp.Get("/list", listCmt)
	cp.Patch("/{id:int}", patchCmt)
	cp.Delete("/{id:int}", delCmt)

	ppp := pa.Party("/preferences")
	mvc.New(ppp).Handle(new(prefCtrl))
	ppp.Get("/get", getPref)
	ppp.Patch("/", patchPref)
}

func submodule(expect1, cond1, expect2, cond2 string) bool {
	return expect1 == cond1 && expect2 == cond2
}

func exists(o interface{}) bool {
	return o != nil
}

func str(content []byte) string {
	return string(content)
}

func content(length int, content []byte) string {
	return truncate(length, string(content))
}

func truncate(length int, s string) string {
	sRune := []rune(s)
	if len(sRune) > length {
		return string(sRune[:length]) + "..."
	}
	return s
}

func date(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func tagSelected(tags []Tag, tag Tag) bool {
	for _, v := range tags {
		if v.ID == tag.ID {
			return true
		}
	}
	return false
}

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
		Name: "tpl/dashboard.html",
		Data: iris.Map{
			"SectionName":    "dashboard",
			"SubSectionName": "dashboard",
			"tab":            dashboard,
			"PostNumber":     pc,
			"CommentNumber":  cc,
			"DaysLastPost":   dp,
			"DayOnline":      do,
		},
	}
}
