package admin

import (
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/services/tags"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type tagCtrl struct {
	Ctx iris.Context
}

func (c *tagCtrl) Get() mvc.View {
	return mvc.View{
		Name: "admin/html/tag.html",
		Data: iris.Map{
			"tab": tag,
		},
	}
}

// GET /admin/tags/list
func listTag(ctx iris.Context) {
	ts, err := TagTable().All()
	if err != nil {
		ts = make([]Tag, 0)
	}
	ctx.JSON(Result(true, "ok", ts))
}

// POST /admin/tags
func createTag(ctx iris.Context) {
	ft := &formTag{}
	err := ctx.ReadJSON(ft)
	if err != nil {
		Logger().Error("tagCtrl", "无法读取数据", Params{}.Err(err))
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "无法读取数据", nil))
		return
	}

	if ft.Name == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "标签名不能为空", nil))
		return
	}

	id, err := TagTable().Create(Params{"name": ft.Name})

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "创建标签失败", nil))
		Logger().Error("tagCtrl", "创建标签失败", Params{"id": id, "name": ft.Name}.Err(err))
		return
	}
	ctx.JSON(Result(true, "ok", Params{"ID": id}))
}

// DELETE /admin/tags/1
func delTag(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "获取不到标签 id", nil))
		return
	}

	if err := TagTable().Delete(id); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "删除标签失败", nil))
		Logger().Error("tagCtrl", "删除标签失败", Params{"id": id}.Err(err))
		return
	}

	ctx.JSON(Result(true, "ok", nil))
}

type formTag struct {
	Name string
}

// PATCH /admin/tags/1
func patchTag(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "获取不到标签 id", nil))
		return
	}

	ft := &formTag{}
	err = ctx.ReadJSON(ft)
	if err != nil {
		Logger().Error("tagCtrl", "无法读取数据", Params{"id": id}.Err(err))
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "无法读取数据", nil))
		return
	}

	if ft.Name == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "标签名不能为空", nil))
		return
	}

	_, err = TagTable().Retrieve(id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "获取不到此标签记录", nil))
		Logger().Error("tagCtrl", "获取不到此标签记录", Params{"id": id}.Err(err))
		return
	}

	if err := TagTable().Update(id, Params{"name": ft.Name}); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "更新标签失败", nil))
		Logger().Error("tagCtrl", "更新标签失败", Params{"id": id}.Err(err))
		return
	}

	ctx.JSON(Result(true, "ok", nil))
}
