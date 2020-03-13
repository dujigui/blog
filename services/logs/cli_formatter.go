package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"sync/atomic"
	"time"
)

type CliFormatter struct {
	c uint32
}

func (f *CliFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	i := atomic.LoadUint32(&f.c)
	atomic.AddUint32(&f.c, 1)

	var b1, b2 []string
	for k, v := range entry.Data {
		if k == "tag" {
			continue
		}
		if v, ok := v.(time.Time); ok {
			b2 = append(b2, fmt.Sprintf("%s=%s", k, v.Format("2006.01.02 15:04:05")))
			continue
		}
		b1 = append(b1, fmt.Sprintf("%s=%s", k, fmt.Sprint(v)))
	}

	// [fatal][tag][001][2006.01.02 15:04:05] message data
	r := fmt.Sprintf("[%03d][%s][%s][%s] %s %s %s\n",
		i,
		strings.ToUpper(entry.Level.String()),
		entry.Data["tag"],
		entry.Time.Format("2006.01.02 15:04:05"),
		entry.Message,
		strings.Join(b1, " "),
		strings.Join(b2, " "),
	)
	return []byte(r), nil
}
