package visitor

import (
	"github.com/dujigui/blog/services/cfg"
	. "github.com/dujigui/blog/services/logs"
	"github.com/dujigui/blog/services/posts"
	. "github.com/dujigui/blog/services/tags"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func Visitor(app *iris.Application) {
	app.HandleDir("/", cfg.Config().GetString("favicon"))

	app.Get("/posts/{id:int}", func(ctx iris.Context) {
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
	})
}

type postsCtrl struct {
	Ctx iris.Context
}

func (c *postsCtrl) GetBy(id int) mvc.View {
	p, err := posts.PostTable().Retrieve(id)
	if err != nil {
		return ErrMsg("无此 ID 文章")
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

	tags, err := TagTable().All()
	if err != nil {
		tags = make([]Tag, 0)
	}

	return mvc.View{
		Name: "tpl/post_editor.html",
		Data: iris.Map{
			"SectionName":    "post",
			"SubSectionName": "editor",
			"Post":           p,
			"Tags":           tags,
		},
	}
}
