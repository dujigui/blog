package admin

import (
	"fmt"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/services/posts"
	. "github.com/dujigui/blog/services/tags"
	. "github.com/dujigui/blog/utils"
	"github.com/iris-contrib/blackfriday"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strings"
)

type postsCtrl struct {
	Ctx iris.Context
}

// GET /admin/posts
func (c *postsCtrl) Get() mvc.View {
	return mvc.View{
		Name: "admin/html/post_list.html",
		Data: iris.Map{
			"tab": postList,
		},
	}
}

// GET /admin/posts/editor
func (c *postsCtrl) GetEditor() mvc.View {
	return mvc.View{
		Name: "admin/html/post_editor.html",
		Data: iris.Map{
			"tab": postEditor,
		},
	}
}

// GET /admin/posts/editor/1
func (c *postsCtrl) GetEditorBy(id int) mvc.View {
	p, err := PostTable().Retrieve(id)
	if err != nil {
		return ErrMsg("无此 ID 文章")
	}

	return mvc.View{
		Name: "admin/html/post_editor.html",
		Data: iris.Map{
			"tab":  postEditor,
			"Post": p,
		},
	}
}

// content 进 content 出
// tagIds 进 tagIds 出
type formPost struct {
	Post
	//TagIDs  string
	Content string
}

// POST /admin/posts
func (c *postsCtrl) Post() interface{} {
	fp := &formPost{}
	if err := c.Ctx.ReadForm(fp); err != nil {
		return ErrMsg("无法读取 form 的数据")
	}
	fp.TagIDs = strings.Join(c.Ctx.FormValues()["TagIDs"], ",")

	p, em := params(fp)
	if em != "" {
		return ErrMsg(em)
	}

	id, err := PostTable().Create(p)
	if err != nil {
		return ErrMsg(fmt.Sprintf("内部错误 %s", err))
	}

	if fp.IsPublished {
		return mvc.Response{Path: fmt.Sprintf("/posts/%d", id)}
	}
	return mvc.Response{Path: "/admin/posts"}
}

// POST /admin/posts/1
func (c *postsCtrl) PostBy(id int) interface{} {
	fp := &formPost{}
	if err := c.Ctx.ReadForm(fp); err != nil {
		return ErrMsg("无法读取 form 的数据")
	}
	fp.TagIDs = strings.Join(c.Ctx.FormValues()["TagIDs"], ",")

	p, em := params(fp)
	if em != "" {
		return ErrMsg(em)
	}

	err := PostTable().Update(id, p)
	if err != nil {
		return ErrMsg(fmt.Sprintf("内部错误 %s", err))
	}

	if fp.IsPublished {
		return mvc.Response{Path: fmt.Sprintf("/posts/%d", id)}
	}
	return mvc.Response{Path: "/admin/posts"}
}

func params(fp *formPost) (Params, string) {
	p := Params{}
	p["type"] = fp.Type
	switch fp.Type {
	case Article:
		if fp.Title == "" {
			return p, "Title 不能为空"
		}
		p["title"] = fp.Title
		if fp.Description == "" {
			return p, "Description 不能为空"
		}
		p["description"] = fp.Description
	case Snippet:
		p["title"] = ""
		if fp.Description == "" {
			return p, "Description 不能为空"
		}
		p["description"] = fp.Description
	case Moment:
		p["title"] = ""
		p["description"] = ""
	default:
		return p, "Type 只能为1、2、3"
	}

	p["cover"] = fp.Cover
	p["is_published"] = fp.IsPublished
	p["content"] = []byte(fp.Content)
	p["tag_ids"] = fp.TagIDs
	return p, ""
}

// GET /admin/posts/list
func listPost(ctx iris.Context) {
	page := ctx.URLParamIntDefault("page", 1)
	limit := ctx.URLParamIntDefault("limit", 10)

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
	ctx.JSON(Result(true, "ok", pp, "page", page, "limit", limit, "total", total))
}

// DELETE /admin/posts/1
func delPost(ctx iris.Context) {
	id := ctx.Params().GetIntDefault("id", -1)
	if id == -1 {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	if err := PostTable().Delete(id); err != nil {
		Logger().Error("postCtrl", "删除文章失败", Params{"id": id}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
	} else {
		ctx.StatusCode(iris.StatusOK)
	}
}

// GET /admin/posts/1
func getPost(ctx iris.Context) {
	id := ctx.Params().GetIntDefault("id", -1)
	if id == -1 {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	p, err := PostTable().Retrieve(id)
	if err != nil {
		Logger().Error("postCtrl", "获取文章失败", Params{"id": id}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, err.Error(), nil))
		return
	}
	pf := formPost{Post: p, Content: string(p.Content)}

	ctx.JSON(Result(true, "ok", pf))
}

func markdown(ctx iris.Context) {
	d, _ := ctx.GetBody()
	ctx.Write(blackfriday.Run(d))
}
