package visitor

import (
	"github.com/dujigui/blog/gateway"
	. "github.com/dujigui/blog/services/comments"
	. "github.com/dujigui/blog/services/posts"
	. "github.com/dujigui/blog/services/tags"
	. "github.com/dujigui/blog/services/users"
	. "github.com/dujigui/blog/utils"
	"github.com/iris-contrib/blackfriday"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"html"
	"html/template"
)

type detailCtrl struct {
	Ctx iris.Context
}

func (c *detailCtrl) GetBy(id int) mvc.View {
	p, err := PostTable().Retrieve(id)

	if err != nil || !p.Publish{
		return ErrMsg("无此 ID 文章")
	}

	ts, err := TagTable().RetrieveIDs(p.TagIDs)
	if err != nil {
		ts = make([]Tag, 0)
	}
	p.Tags = ts

	cs, err := CommentTable().ByPost(p.ID)
	if err != nil {
		cs = make([]Comment, 0)
	}

	for k, v := range cs {
		u, err := UserTable().Retrieve(v.UserID)
		if err != nil {
			cs[k].User = CommentUser{
				ID:       0,
				Avatar:   "/images/avatar.svg",
				Nickname: "unknown",
			}
			continue
		}

		cs[k].User = CommentUser{
			ID:       u.ID,
			Avatar:   u.Avatar,
			Nickname: u.Nickname,
		}
		// 模板引擎会做一次转义
		cs[k].Content = html.UnescapeString(v.Content)
	}

	ok, uid, _ := gateway.Info(c.Ctx)
	var user CommentUser
	if ok && uid != 0 {
		u, err := UserTable().Retrieve(uid)
		if err == nil {
			user.ID = u.ID
			user.Nickname = u.Nickname
			user.Avatar = u.Avatar
		}
	}

	return mvc.View{
		Name: "visitor/html/detail.html",
		Data: iris.Map{
			"tab":       detail,
			"post":      p,
			"comments":  cs,
			"content":   template.HTML(string(blackfriday.Run(p.Content))),
			"needLogin": !ok || uid == 0,
			"user":      user,
		},
	}
}
