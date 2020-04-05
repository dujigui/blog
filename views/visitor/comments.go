package visitor

import (
	"fmt"
	. "github.com/dujigui/blog/services/comments"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"strconv"
)

type formComment struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

func postComment(ctx iris.Context) {
	ok := ctx.Params().GetBoolDefault("ok", false)
	uid := ctx.Params().GetIntDefault("id", 0)
	admin := ctx.Params().GetBoolDefault("admin", false)
	if !ok || uid == 0 {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(Result(false, "请重新登录", nil))
		return
	}

	var fc formComment
	if err := ctx.ReadJSON(&fc); err != nil {
		fmt.Println(err)
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "无法读取 body", nil))
		return
	}

	pid, err := strconv.Atoi(fc.PostID)
	if err != nil || pid <= 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "PostID 不能为空", nil))
		return
	}

	if fc.Content == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "Content 不能为空", nil))
		return
	}

	cid, err := CommentTable().Create(Params{"content": fc.Content, "post_id": fc.PostID, "user_id": uid, "inspect": admin})
	if err != nil || cid <= 0 {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "创建评论失败", nil))
		return
	}

	ctx.JSON(Result(true, "评论成功！审核后显示。", nil))
}
