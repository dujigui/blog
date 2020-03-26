package admin

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
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
	app.HandleDir("/backyard/css", "views/web/admin/css")
	app.HandleDir("/backyard/js", "views/web/admin/js")

	pa := app.Party("/admin").Layout("admin/admin.html")
	mvc.New(pa).Handle(new(adminCtrl))

	pp := pa.Party("/posts")
	mvc.New(pp).Handle(new(postsCtrl))
	pp.Get("/list", listPost)
	pp.Delete("/{id:int}", delPost)
	pp.Get("/{id:int}", getPost)
	pp.Post("/markdown", markdown)

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
