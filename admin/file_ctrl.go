package admin

import (
	"fmt"
	. "github.com/dujigui/blog/services/files"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"time"
)

type filesCtrl struct {
	Ctx iris.Context
}

// GET /admin/files
func (c *filesCtrl) Get() mvc.View {
	return mvc.View{
		Name: "tpl/file.html",
		Data: iris.Map{
			"tab": file,
		},
	}
}

// DELETE /admin/files/1
func delFile(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "获取不到文件 id", nil))
		return
	}

	f, err := FileTable().Retrieve(id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "获取不到此文件记录", nil))
		Logger().Error("fileCtrl", "获取不到此文件记录", Params{"filename": f.Hashed, "id": f.ID}.Err(err))
		return
	}

	if err := FileTable().Delete(id); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, fmt.Sprintf("删除数据库文件记录失败 %s", err), nil))
		Logger().Error("fileCtrl", "删除数据库文件记录失败", Params{"filename": f.Hashed, "id": f.ID}.Err(err))
		return
	}
	if err := Remove(f.Hashed); err != nil {
		Logger().Error("fileCtrl", "删除磁盘文件失败", Params{"filename": f.Hashed, "id": f.ID}.Err(err))
	}

	ctx.JSON(Result(true, "ok", nil))
	Logger().Trace("fileCtrl", "删除文件成功", Params{"filename": f.Hashed, "id": f.ID})
}

// GET /admin/files/list
func listFile(ctx iris.Context) {
	page := ctx.URLParamIntDefault("page", 1)
	limit := ctx.URLParamIntDefault("limit", 10)

	fs, total, err := FileTable().Page(page, limit)
	if err != nil {
		fs = make([]File, 0)
		total = 0
	}
	ctx.JSON(Result(true, "ok", fs, "total", total))
}

// POST /admin/files
func uploadFile(ctx iris.Context) {
	f, info, err := ctx.FormFile("file")
	if err != nil {
		Logger().Error("fileCtrl", "获取文件信息失败", Params{}.Err(err))
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "获取文件信息失败", nil))
		return
	}
	defer f.Close()

	fn := info.Filename
	hashed, err := Save(f, fn)
	if err != nil {
		Logger().Error("fileCtrl", "保存磁盘文件失败", Params{"original": fn, "hashed": hashed}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, fmt.Sprintf("保存磁盘文件失败 %s", err), nil))
		return
	}

	p := Params{"original": fn, "hashed": hashed, "created": time.Now(), "updated": time.Now()}
	id, err := FileTable().Create(p)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, fmt.Sprintf("保存文件数据库记录失败 %s", err), nil))
		Logger().Error("fileCtrl", "保存文件数据库记录失败", Params{"original": fn, "hashed": hashed}.Err(err))
		return
	}
	ff, err := FileTable().Retrieve(id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, err.Error(), nil))
	}

	ctx.JSON(Result(true, "ok", ff))
	return
}
