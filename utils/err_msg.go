package utils

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func ErrMsg(msg string) mvc.View {
	return mvc.View{
		Name:   "tpl/error.html",
		Layout: iris.NoLayout,
		Data: iris.Map{
			"ErrorMessage": msg,
		},
	}
}
