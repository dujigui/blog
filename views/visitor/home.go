package visitor

import (
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/services/posts"
	. "github.com/dujigui/blog/services/tags"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type homeCtrl struct {
	Ctx iris.Context
}

func (c *homeCtrl) Get() mvc.View {
	page := c.Ctx.URLParamIntDefault("page", 1)
	limit := c.Ctx.URLParamIntDefault("limit", 10)

	posts, total, err := PostTable().Page(page, limit)
	if err != nil {
		Logger().Error("postsCtrl", "无法获取文章列表", Params{"total": total, "err": err})
		posts = make([]Post, 0)
		total = 0
	}

	for i, p := range posts {
		if p.TagIDs == "" {
			posts[i].Tags = make([]Tag, 0)
			continue
		}

		tags, err := TagTable().RetrieveIDs(posts[i].TagIDs)
		posts[i].Tags = tags

		if err != nil {
			posts[i].Tags = make([]Tag, 0)
			Logger().Error("postsCtrl", "无法获取文章标签", Params{"post": p.ID, "tag_ids": p.TagIDs}.Err(err))
		}
	}

	type formPost struct {
		Post
		Content string
	}

	var pp []formPost
	for i := range posts {
		s := string(posts[i].Content)

		sRune := []rune(s)
		if len(sRune) > 100 {
			s = string(sRune[:100]) + "..."
		}
		ppp := formPost{Post: posts[i], Content: s}
		pp = append(pp, ppp)
	}

	return mvc.View{
		Name: "visitor/html/home.html",
		Data: iris.Map{
			"tab":     home,
			"posts":   pp,
			"current": page,
			"total":   total/limit+1,
		},
	}
}
