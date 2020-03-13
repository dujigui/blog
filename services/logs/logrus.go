package logs

import (
	"fmt"
	. "github.com/dujigui/blog/services/cfg"
	. "github.com/dujigui/blog/utils"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
	"time"
)

var l = logger{}

type logger struct {
}

func init() {
	if Config().GetBool("production") {
		updateOutput()
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetFormatter(&CliFormatter{})
		logrus.SetOutput(os.Stdout)
		logrus.SetLevel(logrus.TraceLevel)
	}
}

func updateOutput() {
	rotate := Config().GetString("logs.rotate")
	switch rotate {
	case "hour":
		updateFile(time.Hour)
	case "day":
		updateFile(time.Hour * 24)
	default:
		log.Fatalf("参数错误 logs.rotate=%s", rotate)
	}
}

func updateFile(d time.Duration) {
	dir := Config().GetString("logs.dir")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("创建日志目录失败 dir=%s", dir)
	}

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	fp := fileName(now, d == time.Hour*24)
	f := setFile(dir, fp)

	go func() {
		for {
			next := now.Add(d)
			t := time.NewTimer(next.Sub(now))
			<-t.C
			t.Stop()
			fp = fileName(now, d != time.Hour)
			tf := setFile(dir, fp)
			if err := f.Close(); err != nil {
				log.Printf("关闭日志文件失败 err=%s old=%s\n", err, fp)
			}
			f = tf
			now = next
		}
	}()
}

func fileName(t time.Time, daily bool) string {
	if daily {
		return fmt.Sprintf("%d-%d-%d.log", t.Year(), t.Month(), t.Day())
	} else {
		return fmt.Sprintf("%d-%d-%d-%d.log", t.Year(), t.Month(), t.Day(), t.Hour())
	}
}

func setFile(dir, fp string) *os.File {
	// create if not exists
	// append if exists
	f, err := os.OpenFile(filepath.Join(dir, fp), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("打开日志文件失败 err=%s fileName=%s", err, fp)
	}
	logrus.SetOutput(f)
	return f
}

func (l *logger) Trace(tag, msg string, params Params) {
	var entry = logrus.WithField(TAG, tag)
	for k, v := range params {
		entry = entry.WithField(k, v)
	}
	entry.Trace(msg)
}

func (l *logger) Debug(tag, msg string, params Params) {
	var entry = logrus.WithField(TAG, tag)
	for k, v := range params {
		entry = entry.WithField(k, v)
	}
	entry.Debug(msg)
}

func (l *logger) Info(tag, msg string, params Params) {
	var entry = logrus.WithField(TAG, tag)
	for k, v := range params {
		entry = entry.WithField(k, v)
	}
	entry.Info(msg)
}

func (l *logger) Warning(tag, msg string, params Params) {
	var entry = logrus.WithField(TAG, tag)
	for k, v := range params {
		entry = entry.WithField(k, v)
	}
	entry.Warning(msg)
}

func (l *logger) Error(tag, msg string, params Params) {
	var entry = logrus.WithField(TAG, tag)
	for k, v := range params {
		entry = entry.WithField(k, v)
	}
	entry.Error(msg)
}

func (l *logger) Panic(tag, msg string, params Params) {
	var entry = logrus.WithField(TAG, tag)
	for k, v := range params {
		entry = entry.WithField(k, v)
	}
	entry.Panic(msg)
}

func (l *logger) Fatal(tag, msg string, params Params) {
	var entry = logrus.WithField(TAG, tag)
	for k, v := range params {
		entry = entry.WithField(k, v)
	}
	entry.Fatal(msg)
}
