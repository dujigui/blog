package logs

import (
	"fmt"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12/context"
	ml "github.com/kataras/iris/v12/middleware/logger"
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
	return ml.New(
		ml.Config{
			Status:  true,
			IP:      true,
			Method:  true,
			Path:    true,
			Query:   true,
			Columns: true,
			LogFunc: LogReq,
		})
}

func LogReq(endTime time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
	p := Params{
		"endTime": fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", endTime.Year(), endTime.Month(), endTime.Day(), endTime.Hour(), endTime.Minute(), endTime.Second()),
		"latency": latency,
		"status":  status,
		"ip":      ip,
		"method":  method,
		"path":    path,
	}
	if headerMessage != nil {
		p["headerMessage"] = headerMessage
	}
	if message != nil {
		p["message"] = message
	}

	Logger().Debug("logReq", "", p)
}
