package visitor

import (
	"fmt"
	. "github.com/dujigui/blog/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type aboutCtrl struct {
	Ctx iris.Context
}

func (c *aboutCtrl) Get() interface{} {
	apid := Pref().AboutPostID
	if apid == 0 {
		return mvc.Response{Path: "/"}
	}

	return mvc.Response{Path: fmt.Sprintf("/posts/%d", apid)}
}
