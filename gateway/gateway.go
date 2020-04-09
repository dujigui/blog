package gateway

import (
	"fmt"
	. "github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/users"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	letters  = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	prefix   = randSeq(10) + "_"
	keyOk    = prefix + "ok"
	keyID    = prefix + "id"
	keyAdmin = prefix + "admin"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	rand.Shuffle(len(b), func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})
	return string(b)
}

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
	ctx.Params().Set(keyOk, strconv.FormatBool(ok))
	if ok {
		ctx.Params().Set(keyID, strconv.Itoa(id))
		ctx.Params().Set(keyAdmin, strconv.FormatBool(admin))
	}

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
	ok = ctx.Params().GetBoolDefault(keyOk, false)
	uid = ctx.Params().GetIntDefault(keyID, 0)
	admin = ctx.Params().GetBoolDefault(keyAdmin, false)
	return
}
