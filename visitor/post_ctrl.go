package visitor

import (
	. "github.com/dujigui/blog/services/logs"
	"github.com/dujigui/blog/services/posts"
	. "github.com/dujigui/blog/services/tags"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
)

type postCtrl struct {
	Ctx iris.Context
}

//func (c *postCtrl) Get() mvc.View {
//
//}


func GetBy(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.Text("no id")
		return
	}
	p, err := posts.PostTable().Retrieve(id)
	if err != nil {
		ctx.Text("no post")
	}

	if p.TagIDs == "" {
		p.Tags = make([]Tag, 0)
	} else {
		p.Tags, err = TagTable().RetrieveIDs(p.TagIDs)
		if err != nil {
			p.Tags = make([]Tag, 0)
			Logger().Error("postsCtrl", "无法获取文章标签", Params{"post": p.ID, "tag_ids": p.TagIDs}.Err(err))
		}
	}

	ctx.JSON(p)
	return
}