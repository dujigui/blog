package visitor

import (
	. "github.com/dujigui/blog/services/posts"
	. "github.com/dujigui/blog/services/tags"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"sort"
	"strconv"
	"strings"
)

type searchCtrl struct {
	Ctx iris.Context
}

type tplPost struct {
	Post
	Content string
}

type tplTag struct {
	Tag
	Posts []tplPost
}

//todo 使用 redis 缓存数据
func (c *searchCtrl) Get() mvc.View {
	ts, err := TagTable().All()
	if err != nil {
		ts = make([]Tag, 0)
	}

	ps, err := PostTable().All()
	if err != nil {
		ps = make([]Post, 0)
	}

	tpm := make(map[int][]tplPost)
	for pk, pv := range ps {
		tids := strings.Split(pv.TagIDs, ",")
		for _, tid := range tids {
			id, err := strconv.Atoi(tid)
			if err == nil {
				tp := tplPost{Post: ps[pk]}
				if tp.Type != Moment {
					tp.Content = ""
				} else {
					tp.Content = string(tp.Post.Content)
				}
				tpm[id] = append(tpm[id], tp)
			}
		}
	}

	var tts []tplTag
	for _, tv := range ts {
		tts = append(tts, tplTag{Tag: tv, Posts: tpm[tv.ID]})
	}
	sort.SliceStable(tts, func(i, j int) bool {
		return strings.Compare(tts[i].Name, tts[j].Name) < 0
	})

	return mvc.View{
		Name: "visitor/html/search.html",
		Data: iris.Map{
			"tab":  search,
			"tags": tts,
		},
	}
}
