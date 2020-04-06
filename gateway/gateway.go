package gateway

import (
	"fmt"
	. "github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/users"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/url"
	"strconv"
	"strings"
)

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

	if Pref().Init <= 0 && ctx.Path() != "/init" {
		ctx.Redirect("/init")
		return
	}

	ok, id, admin := ParseToken(ctx.GetCookie("token"))
	if ok {
		ctx.Params().Set("__id", strconv.Itoa(id))
		ctx.Params().Set("__admin", strconv.FormatBool(admin))
	}
	ctx.Params().Set("__ok", strconv.FormatBool(ok))

	if strings.HasPrefix(ctx.Path(), "/admin") {
		if !ok || !admin {
			ctx.Redirect(fmt.Sprintf("/login?redirect=%s", url.QueryEscape(ctx.Path())))
			return
		}
	}

	if strings.HasPrefix(ctx.Path(), "/comments") {
		if !ok {
			ctx.Redirect(fmt.Sprintf("/login?redirect=%s", url.QueryEscape(ctx.Path())))
			return
		}
	}
	ctx.Next()
}

func Info(ctx iris.Context) (ok bool, uid int, admin bool) {
	ok = ctx.Params().GetBoolDefault("__ok", false)
	uid = ctx.Params().GetIntDefault("__id", 0)
	admin = ctx.Params().GetBoolDefault("__admin", false)
	return
}
