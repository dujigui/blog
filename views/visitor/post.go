package visitor

import (
	. "github.com/dujigui/blog/services/comments"
	. "github.com/dujigui/blog/services/posts"
	. "github.com/dujigui/blog/services/tags"
	. "github.com/dujigui/blog/utils"
	"github.com/iris-contrib/blackfriday"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"html/template"
)

type detailCtrl struct {
	Ctx iris.Context
}

func (c *detailCtrl) GetBy(id int) mvc.View {
	p, err := PostTable().Retrieve(id)

	if err != nil {
		return ErrMsg("无此 ID 文章")
	}

	ts, err := TagTable().RetrieveIDs(p.TagIDs)
	if err != nil {
		ts = make([]Tag, 0)
	}
	p.Tags = ts

	cs, err := CommentTable().RetrieveByPost(p.ID)
	if err != nil {
		cs = make([]Comment, 0)
	}

	return mvc.View{
		Name: "visitor/html/detail.html",
		Data: iris.Map{
			"tab":      detail,
			"post":     p,
			"comments": cs,
			"content":  template.HTML(string(blackfriday.Run(p.Content))),
		},
	}
}
