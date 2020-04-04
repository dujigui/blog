package logs

import (
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"strconv"
	"time"
)

const (
	TAG = "tag"
)

func Logger() Log {
	return &l
}

type Log interface {
	Trace(tag, msg string, params Params)
	Debug(tag, msg string, params Params)
	Info(tag, msg string, params Params)
	Warning(tag, msg string, params Params)
	Error(tag, msg string, params Params)
	Panic(tag, msg string, params Params)
	Fatal(tag, msg string, params Params)
}

func ReqLogger() context.Handler {
	return func(ctx iris.Context) {
		start := time.Now()
		ctx.Next()
		end := time.Now()
		latency := end.Sub(start)
		status := strconv.Itoa(ctx.GetStatusCode())
		ip := ctx.RemoteAddr()
		method := ctx.Method()
		path := ctx.Request().RequestURI
		logReq(latency, status, ip, method, path)
	}
}

func logReq(latency time.Duration, status, ip, method, path string) {
	p := Params{
		"reqLatency": latency,
		"reqStatus":  status,
		"reqIP":      ip,
		"reqMethod":  method,
		"reqPath":    path,
	}
	Logger().Debug("logReq", "", p)
}
