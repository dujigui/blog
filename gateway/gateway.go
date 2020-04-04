package gateway

import (
	. "github.com/dujigui/blog/services"
	"github.com/kataras/iris/v12/context"
	"strings"
)

const (
	ADMIN = "/admin"
)

// 401 unauthorized
// 403 forbidden
// https://docs.iris-go.com/iris/request-authentication
func Gateway(ctx context.Context) {
	if strings.HasPrefix(ctx.Path(), "/files") ||
		strings.HasPrefix(ctx.Path(), "/layui") ||
		strings.HasPrefix(ctx.Path(), "/prism") ||
		strings.HasPrefix(ctx.Path(), "/images") ||
		strings.HasPrefix(ctx.Path(), "/backyard/js") ||
		strings.HasPrefix(ctx.Path(), "/backyard/css") ||
		strings.HasPrefix(ctx.Path(), "/visitor/css") ||
		strings.HasPrefix(ctx.Path(), "/visitor/js") {
		ctx.Next()
		return
	}

	if Pref().Init == 0 && ctx.Path() != "/init" {
		ctx.Redirect("/init")
		return
	}
	if strings.HasPrefix(ctx.Path(), "/admin")  {
		if err := auth(ctx); err != nil {
			ctx.Redirect("/login")
			return
		}
	}
	// todo 添加用户评论权限检查

	ctx.Next()
}

func auth(ctx context.Context) error {
	return nil
}
