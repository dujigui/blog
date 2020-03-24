package admin

import (
	. "github.com/dujigui/blog/services/comments"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type cmtCtrl struct {
	Ctx iris.Context
}

// GET /admin/comments
func (c *cmtCtrl) Get() mvc.View {
	return mvc.View{
		Name: "tpl/comment.html",
		Data: iris.Map{
			"tab": comment,
		},
	}
}

// POST /admin/comments
func patchCmt(ctx iris.Context) {
	id := ctx.Params().GetIntDefault("id", -1)
	if id == -1 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "无法读取id", nil))
		return
	}

	var c = &Comment{}

	if err := ctx.ReadJSON(c); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "无法读取body", nil))
		return
	}

	p := Params{"content": c.Content}

	err := CommentTable().Update(id, p)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "更新评论失败", nil))
		Logger().Error("tagCtrl", "更新评论失败", Params{"id": id}.Err(err))
		return
	}

	cc, err := CommentTable().Retrieve(id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "获取新评论失败", nil))
		Logger().Error("tagCtrl", "获取新评论失败", Params{"id": id}.Err(err))
		return
	}
	ctx.JSON(Result(true, "ok", cc))
}

// GET /admin/comments/list
func listCmt(ctx iris.Context) {
	page := ctx.URLParamIntDefault("page", 1)
	limit := ctx.URLParamIntDefault("limit", 10)

	cmts, total, err := CommentTable().Page(page, limit)
	if err != nil {
		Logger().Error("cmtCtrl", "无法获取评论列表", Params{"total": total, "err": err})
		cmts = make([]Comment, 0)
		total = 0
	}

	ctx.JSON(Result(true, "ok", cmts, "page", page, "limit", limit, "total", total))
}

// DELETE /admin/comments/1
func delCmt(ctx iris.Context) {
	id := ctx.Params().GetIntDefault("id", -1)
	if id == -1 {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	if err := CommentTable().Delete(id); err != nil {
		Logger().Error("cmtCtrl", "删除评论失败", Params{"id": id}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
	} else {
		ctx.StatusCode(iris.StatusOK)
	}
}